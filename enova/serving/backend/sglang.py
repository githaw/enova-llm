import dataclasses
from enova.common.logger import LOGGER
from enova.common.config import CONFIG
from enova.serving.backend.base import BaseBackend


@dataclasses.dataclass
class SglangBackend(BaseBackend):
    def __post_init__(self):
        """Initialize the SglangBackend specific components."""
        self._create_app()

    def _create_app(self):
        from sglang.srt.server import app as sglang_app, launch_engine
        from sglang.srt.server_args import ServerArgs
        from sglang.srt.utils import add_prometheus_middleware
        from sglang.srt.metrics.func_timer import enable_func_timer

        if not hasattr(self, 'model'):
            raise RuntimeError("Model path must be specified")

        if 'tensor_parallel_size' in CONFIG.sglang:
            CONFIG.sglang['tp_size'] = CONFIG.sglang.pop('tensor_parallel_size')
        server_args = ServerArgs(model_path=self.model, **CONFIG.sglang)
        launch_engine(server_args)
        add_prometheus_middleware(sglang_app)
        enable_func_timer()

        self.app = sglang_app

        @self.app.get("/v1/model/info/args")
        async def get_engine_args():
            return {"code": 0, "result": server_args}

        LOGGER.info("SGLangBackend FastAPI app created and routes defined.")
