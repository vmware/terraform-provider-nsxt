#!/bin/bash -x

this_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
provider_file=${this_dir}/../nsxt/provider.go

# Add 1 to include order comment
datasource_start_line=$(expr $(grep -n DataSourcesMap ${provider_file} | cut -f1 -d:) + 1)
# Add 1 to include order comment
resource_start_line=$(expr $(grep -n ResourcesMap ${provider_file} | cut -f1 -d:) + 1)
total_lines=$(wc -l ${provider_file} | awk '{print $1}')

tmp_fname=$(mktemp -u ${TMPDIR}tmp_provider.XXXXX)

head -n ${datasource_start_line} ${provider_file} > ${tmp_fname}

datasource_lines=$(tail -n $(expr ${total_lines} - ${datasource_start_line} - 1) ${provider_file} | grep -n '},' | head -n1 | cut -f1 -d:)
tail -n $(expr ${total_lines} - ${datasource_start_line}) ${provider_file} | head -n ${datasource_lines} | sort >> ${tmp_fname}
tail -n $(expr ${total_lines} - ${datasource_start_line} - ${datasource_lines}) ${provider_file} | head -n $(expr ${resource_start_line} - ${datasource_start_line} - ${datasource_lines}) >> ${tmp_fname}

resource_lines=$(tail -n $(expr ${total_lines} - ${resource_start_line} - 1) ${provider_file} | grep -n '},' | head -n1 | cut -f1 -d:)
tail -n $(expr ${total_lines} - ${resource_start_line}) ${provider_file} | head -n ${resource_lines} | sort >> ${tmp_fname}
tail -n $(expr ${total_lines} - ${resource_start_line} - ${resource_lines}) ${provider_file}  >> ${tmp_fname}

cmp -s ${tmp_fname} ${provider_file}
result=$?
if [ ${result} -ne 0 ]; then
  echo "ResourcesMap or DataSourcesMap aren't properly sorted"
else
  echo "Check successful!"
fi
rm -rf ${tmp_fname}

exit ${result}