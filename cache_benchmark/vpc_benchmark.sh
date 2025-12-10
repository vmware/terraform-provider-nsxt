#!/bin/bash

set -e

PLAN="$1"
#enable cache for the run
export NSXT_ENABLE_CACHE="$2"

VPC_VALUES=(5)
SUBNET_VALUES=(300)
OUT_FILE="vpc_bench_results.txt"


# Header (only once)
if [ ! -f "$OUT_FILE" ]; then
  printf "%-20s %-10s %-15s %-15s %-15s %-10s %-10s %-15s\n" \
  "TIMESTAMP" "VPC_COUNT" "SUBNET_COUNT" "PLAN_SECONDS" "APPLY_SECONDS" "STATUS" "PLAN" "CACHE_ENABLED" >> "$OUT_FILE"
fi

for vpc in "${VPC_VALUES[@]}"; do
  for subnet in "${SUBNET_VALUES[@]}"; do

    PREFIX="run-vpc${vpc}-sub${subnet}" # prefix for each run

    echo "============================================================="
    echo "Running test: prefix=$PREFIX vpc_count=$vpc subnet_count=$subnet"
    echo "============================================================="

    STATUS="success"

    # -------------------------- PLAN --------------------------
    echo "Running PLAN..."

    PLAN_START=$(date +%s)
    if ! terraform plan -no-color \
        -var vpc_prefix="$PREFIX" \
        -var vpc_count="$vpc" \
        -var subnet_count="$subnet" >/dev/null 2>&1; then
      STATUS="plan_failed"
    fi
    PLAN_END=$(date +%s)
    PLAN_DURATION=$((PLAN_END - PLAN_START))

    # -------------------------- APPLY -------------------------
    echo "Running APPLY..."

    APPLY_START=$(date +%s)
    if ! terraform apply -auto-approve -no-color \
        -var vpc_prefix="$PREFIX" \
        -var vpc_count="$vpc" \
        -var subnet_count="$subnet" >/dev/null 2>&1; then
      STATUS="apply_failed"
    fi
    APPLY_END=$(date +%s)
    APPLY_DURATION=$((APPLY_END - APPLY_START))

    # # ------------------------ DESTROY -------------------------
    # echo "sleep 5s before running the destroy"
    # sleep 5

    # echo "Running DESTROY..."
    # terraform destroy -auto-approve \
    #   -var vpc_prefix="$PREFIX" \
    #   -var vpc_count="$vpc" \
    #   -var subnet_count="$subnet" >/dev/null 2>&1 || true

    # # ------------------------ LOG RESULTS ----------------------
    TIMESTAMP=$(date "+%Y-%m-%d %H:%M:%S")

    printf "%-20s %-10s %-15s %-15s %-15s %-10s %-10s %-15s\n" \
      "$TIMESTAMP" "$vpc" "$subnet" "$PLAN_DURATION" "$APPLY_DURATION" "$STATUS" "$PLAN" "$NSXT_ENABLE_CACHE" >> "$OUT_FILE"

    echo "Sleeping for 10 seconds..."
    sleep 10

  done
done

echo "Done! Results written to $OUT_FILE"