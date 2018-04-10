#!/bin/sh

SCRIPT_NAME=$0
SCRIPT_DIR=${SCRIPT_NAME%/*}
cd ${SCRIPT_DIR}

APP_LIST="gateway server"
PACKAGE="http"

ERROR_FILE=${SCRIPT_DIR}/errorcode.md
OUT_FILE=${SCRIPT_DIR}/${PACKAGE}errcode.go

if [[ ! -f ${ERROR_FILE} ]];then
  continue
fi

cat << EOF > ${OUT_FILE}
package ${PACKAGE}

const(
   OK = 0
EOF

cat ${ERROR_FILE} | awk -F'-' '{
      if(NF==3){
         printf("ERR_%s = %d\n", $2, $1);
      }else if($0 ~ /^#[ ]*/){
         sub("#", "//", $0);
         print $0
      }
    }' >> ${OUT_FILE}


cat << EOF >> ${OUT_FILE}
)

var ErrDesc = map[uint]string{
EOF

cat ${ERROR_FILE} | awk -F'-' '{
      if(NF==3){
         printf("ERR_%s : \"%s\",\n", $2, $3);
      }else if($0 ~ /^#[ ]*/){
         sub("#", "//", $0);
         print $0
      }
    }' >> ${OUT_FILE}

echo "}" >> ${OUT_FILE}

gofmt -w ${OUT_FILE}
cat ${OUT_FILE}
