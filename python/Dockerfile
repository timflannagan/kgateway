FROM --platform=$TARGETPLATFORM python:3.12-slim-bookworm AS builder-image

ARG DEBIAN_FRONTEND=noninteractive
ARG TARGETPLATFORM

# Ensure proper directory structure for apt caching
RUN rm -f /etc/apt/apt.conf.d/docker-clean && \
    mkdir -p /var/cache/apt/archives/partial

RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
    --mount=type=cache,target=/var/lib/apt,sharing=locked \
    apt-get update && \
    apt-get install --no-install-recommends -y \
        build-essential \
        libpq-dev && \
    rm -rf /var/lib/apt/lists/*

# Setup virtual environment and PATH in one layer
RUN python3 -m venv /opt/venv && \
    /opt/venv/bin/pip install --no-cache-dir --upgrade pip wheel

ENV VIRTUAL_ENV=/opt/venv \
    PATH="/opt/venv/bin:$PATH" \
    TIKTOKEN_CACHE_DIR="/opt/tiktoken" \
    TLDEXTRACT_CACHE="/opt/tldextract"

ARG PYTHON_DIR
COPY $PYTHON_DIR/requirements.txt .

# Install dependencies with build caching
RUN --mount=type=cache,target=/root/.cache \
    pip install --no-cache-dir --prefer-binary -r requirements.txt && \
    # Pre-download models in the same layer
    python3 -c "import tiktoken; tiktoken.get_encoding('cl100k_base')" && \
    python3 -m tldextract --update

# Runner Stage
FROM --platform=$TARGETPLATFORM python:3.12-slim-bookworm AS runner-image

ARG PYTHON_DIR

# Create user and directory structure in one layer
RUN useradd --create-home kgateway && \
    mkdir -p /home/kgateway/python

# Copy artifacts from builder
COPY --from=builder-image /opt/venv /opt/venv
COPY --from=builder-image /opt/tldextract /opt/tldextract
COPY --from=builder-image /opt/tiktoken /opt/tiktoken
COPY $PYTHON_DIR/ai_extension /home/kgateway/python

# Set permissions in one layer
RUN chown -R kgateway:kgateway /opt/tldextract /opt/tiktoken /home/kgateway && \
    chmod 755 /opt/tldextract /opt/tiktoken

# Switch to non-root user
USER kgateway
WORKDIR /home/kgateway

# Set all environment variables in one layer
ENV PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1 \
    VIRTUAL_ENV=/opt/venv \
    PATH="/opt/venv/bin:$PATH" \
    TLDEXTRACT_CACHE="/opt/tldextract" \
    TIKTOKEN_CACHE_DIR="/opt/tiktoken"

CMD ["python3", "python/main.py"]
