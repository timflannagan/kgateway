FROM --platform=$TARGETPLATFORM python:3.12-slim-bookworm AS builder-image

ARG DEBIAN_FRONTEND=noninteractive
ARG TARGETPLATFORM

RUN --mount=type=cache,target=/var/cache/apt \
    apt-get update && apt-get install --no-install-recommends -y build-essential libpq-dev && \
    rm -rf /var/lib/apt/lists/*

RUN python3 -m venv /opt/venv
ENV VIRTUAL_ENV=/opt/venv
ENV PATH="/opt/venv/bin:$PATH"

ARG PYTHON_DIR

COPY $PYTHON_DIR/requirements.txt .
RUN pip3 install --no-cache-dir wheel
RUN --mount=type=cache,target=/root/.cache/pip \
    pip3 install --no-cache-dir --prefer-binary -r requirements.txt

# Pre-download models and update tldextract
ENV TIKTOKEN_CACHE_DIR="/opt/tiktoken"
RUN python3 -c "import tiktoken; tiktoken.get_encoding('cl100k_base')"
ENV TLDEXTRACT_CACHE="/opt/tldextract"
RUN python3 -m tldextract --update

# Runner Stage
FROM --platform=$TARGETPLATFORM python:3.12-slim-bookworm AS runner-image

ARG PYTHON_DIR

RUN useradd --create-home kgateway

COPY --from=builder-image /opt/venv /opt/venv
COPY --from=builder-image /opt/tldextract /opt/tldextract
COPY --from=builder-image /opt/tiktoken /opt/tiktoken

RUN chown -R kgateway:kgateway /opt/tldextract /opt/tiktoken && \
    chmod 755 /opt/tldextract /opt/tiktoken

USER kgateway
RUN mkdir -p /home/kgateway/python
WORKDIR /home/kgateway
COPY $PYTHON_DIR/ai_extension /home/kgateway/python

# Environment configurations for runtime
ENV PYTHONUNBUFFERED=1
ENV PYTHONDONTWRITEBYTECODE=1
ENV VIRTUAL_ENV=/opt/venv
ENV PATH="/opt/venv/bin:$PATH"
ENV TLDEXTRACT_CACHE="/opt/tldextract"
ENV TIKTOKEN_CACHE_DIR="/opt/tiktoken"

CMD ["python3", "python/main.py"]
