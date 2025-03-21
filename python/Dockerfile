FROM --platform=$TARGETPLATFORM python:3.12-slim-bookworm AS builder-image

ARG DEBIAN_FRONTEND=noninteractive
ARG TARGETPLATFORM

RUN rm -f /etc/apt/apt.conf.d/docker-clean && \
    mkdir -p /var/cache/apt/archives/partial

RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
    --mount=type=cache,target=/var/lib/apt,sharing=locked \
    apt-get update && \
    apt-get install --no-install-recommends -y \
        build-essential=12.9 \
        libpq-dev=15.* && \
    rm -rf /var/lib/apt/lists/*

RUN python3 -m venv /opt/venv && \
    /opt/venv/bin/pip install --no-cache-dir --upgrade pip wheel setuptools

ENV VIRTUAL_ENV=/opt/venv \
    PATH="/opt/venv/bin:$PATH" \
    TIKTOKEN_CACHE_DIR="/opt/tiktoken" \
    TLDEXTRACT_CACHE="/opt/tldextract"

ARG PYTHON_DIR

COPY $PYTHON_DIR/requirements.txt /tmp/requirements.txt
RUN --mount=type=cache,target=/root/.cache/pip,sharing=locked \
    pip install --no-cache-dir --prefer-binary -r /tmp/requirements.txt

RUN --mount=type=cache,target=/root/.cache,sharing=locked \
    python3 -c "import tiktoken; tiktoken.get_encoding('cl100k_base')" && \
    python3 -m tldextract --update

FROM --platform=$TARGETPLATFORM python:3.12-slim-bookworm AS runner-image

ARG PYTHON_DIR

RUN useradd --create-home kgateway && \
    mkdir -p /home/kgateway/python

COPY --from=builder-image /opt/venv /opt/venv
COPY --from=builder-image /opt/tldextract /opt/tldextract
COPY --from=builder-image /opt/tiktoken /opt/tiktoken
COPY $PYTHON_DIR/ai_extension /home/kgateway/python

RUN chown -R kgateway:kgateway /opt/tldextract /opt/tiktoken /home/kgateway && \
    chmod 755 /opt/tldextract /opt/tiktoken

USER kgateway
WORKDIR /home/kgateway

ENV PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1 \
    VIRTUAL_ENV=/opt/venv \
    PATH="/opt/venv/bin:$PATH" \
    TLDEXTRACT_CACHE="/opt/tldextract" \
    TIKTOKEN_CACHE_DIR="/opt/tiktoken"

CMD ["python3", "python/main.py"]
