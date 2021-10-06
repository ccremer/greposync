# These are some common variables for Make

IMG_TAG ?= latest

# Image URL to use all building/pushing image targets
CONTAINER_IMG ?= ghcr.io/ccremer/greposync:$(IMG_TAG)

GOASCIIDOC_OUT_ROOT ?= docs/modules/developer-guide
GOASCIIDOC_OUT_GODOC_PATH ?= $(GOASCIIDOC_OUT_ROOT)/pages/ref-domain.adoc

GOASCIIDOC_ARGS_BUILD ?= -o $(GOASCIIDOC_OUT_GODOC_PATH) --templatedir docs/config-templates
GOASCIIDOC_CMD = go run github.com/mariotoffia/goasciidoc $(GOASCIIDOC_ARGS_BUILD)

REF_CONFIG_PATH ?= docs/modules/ROOT/examples/config.yaml
REF_LABELS_PATH ?= docs/modules/ROOT/examples/labels.yaml
