#!/bin/sh


echo "SEWOBE Anfragen laufen"

curl -X POST -F "USERNAME=$SEWOBEUSER" -F "PASSWORT=$SEWOBEPASSWORD" -F 'AUSWERTUNG_ID=70' $SEWOBEURL > "/go/src/app/testdata/sewobe27.json"
curl -X POST -F "USERNAME=$SEWOBEUSER" -F "PASSWORT=$SEWOBEPASSWORD" -F 'AUSWERTUNG_ID=102' $SEWOBEURL > "/go/src/app/testdata/sewobe102.json"
curl -X POST -F "USERNAME=$SEWOBEUSER" -F "PASSWORT=$SEWOBEPASSWORD" -F 'AUSWERTUNG_ID=158' $SEWOBEURL > "/go/src/app/testdata/sewobe158.json"


echo "Ende von Anfragen"
