# SPDX-FileCopyrightText: Copyright Contributors to the Gardener project
#
# SPDX-License-Identifier: Apache-2.0

TOOLS_DIR     := hack/tools
TOOLS_BIN_DIR := $(TOOLS_DIR)/bin

# Tool Binaries
GOSEC ?= $(TOOLS_BIN_DIR)/gosec

# Tool Versions
GOSEC_VERSION ?= v2.25.0

$(GOSEC):
	@GOSEC_VERSION=$(GOSEC_VERSION) $(TOOLS_DIR)/install-gosec.sh
