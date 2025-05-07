import asyncio
import dataclasses
from enova.common.logger import LOGGER
from enova.common.config import CONFIG
from enova.common.constant import VllmMode
from enova.serving.backend.base import BaseBackend


class CustomDict(dict):

    def __getattribute__(self, name: str):
        if name in self:
            return self[name]
        return None


@dataclasses.dataclass
class VllmBackend(BaseBackend):
    def __post_init__(self):
        """"""

    def _create_app(self):
        vllm_mode = CONFIG.vllm.pop("vllm_mode", VllmMode.NORMAL.value)
        from vllm.engine.arg_utils import AsyncEngineArgs
        from vllm.engine.async_llm_engine import AsyncLLMEngine
        from vllm.transformers_utils.tokenizer import get_tokenizer
        import torch

        if not torch.cuda.is_available():
            raise RuntimeError("vLLM Photon requires CUDA runtime")

        if vllm_mode == VllmMode.NORMAL.value:
            from vllm.entrypoints import api_server

            engine_args = AsyncEngineArgs(model=self.model, **CONFIG.vllm)
            engine = AsyncLLMEngine.from_engine_args(engine_args)
            engine_model_config = asyncio.run(engine.get_model_config())
            max_model_len = engine_model_config.max_model_len

            api_server.served_model = self.model
            api_server.engine = engine
            api_server.max_model_len = max_model_len
            api_server.tokenizer = get_tokenizer(
                engine_args.tokenizer,
                tokenizer_mode=engine_args.tokenizer_mode,
                trust_remote_code=engine_args.trust_remote_code,
            )
        elif vllm_mode == VllmMode.OPENAI.value:
            from vllm.entrypoints.openai import api_server

            class CustomArgParser(api_server.FlexibleArgumentParser):

                def parse_args(self, args=None, namespace=None):
                    args, _ = self.parse_known_args(args, namespace)
                    return args

            parser = CustomArgParser(description="vLLM OpenAI-Compatible RESTful API server.")
            parser = api_server.make_arg_parser(parser)
            args = parser.parse_args()

            current_engine_args = {k: v for k, v in CONFIG.vllm.items() if k in AsyncEngineArgs.__dataclass_fields__}
            engine_args = AsyncEngineArgs(model=self.model, **current_engine_args)
            engine = AsyncLLMEngine.from_engine_args(engine_args, usage_context=api_server.UsageContext.OPENAI_API_SERVER)
            engine_model_config = asyncio.run(engine.get_vllm_config())
            api_server.engine = engine
            api_server.async_engine_client = engine
            api_server.engine_args = engine_args
            api_server.app = api_server.build_app(args)
            asyncio.run(api_server.init_app_state(api_server.async_engine_client, engine_model_config, api_server.app.state, args))

        else:
            raise ValueError(f"vllm_mode: {vllm_mode} is not support")
        LOGGER.info(f"CONFIG.vllm: {CONFIG.vllm}")

        self.app = api_server.app
        cur_app = api_server.app

        @cur_app.get("/v1/model/info/args")
        async def get_engine_args():
            return {"code": 0, "result": engine_args}
