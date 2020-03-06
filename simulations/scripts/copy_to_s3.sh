FILEPATH={{ .Target }}
FILENAME={{ .Filename }}
if [ -f "${FILEPATH}" ]; then
  ln -s ${FILEPATH} /tmp/${FILENAME}
else
  cd ${FILEPATH}
  tar czf /tmp/${FILENAME} *
fi
aws s3 cp /tmp/${FILENAME} {{ .Bucket }}
