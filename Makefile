.PHONY: download
download:
	mkdir -p data

	# Collocations: https://www.clarin.si/repository/xmlui/handle/11356/1250
	cd data && wget -c https://www.clarin.si/repository/xmlui/bitstream/handle/11356/1250/KSSS.zip

	# Thesaurus: https://www.clarin.si/repository/xmlui/handle/11356/1166
	cd data && wget -c https://www.clarin.si/repository/xmlui/bitstream/handle/11356/1166/CJVT_Thesaurus-v1.0.zip

	# Slolex: https://www.clarin.si/repository/xmlui/handle/11356/1230
	#cd data && wget -c https://www.clarin.si/repository/xmlui/handle/11356/1230/allzip -O slolex.zip
	cd data && wget -c https://www.clarin.si/repository/xmlui/bitstream/handle/11356/1230/Sloleks2.0.LMF.zip
	cd data && wget -c https://www.clarin.si/repository/xmlui/bitstream/handle/11356/1230/Sloleks2.0.MTE.zip
	cd data && unzip Sloleks2.0.MTE.zip

	cd data && wget -c https://www.clarin.si/repository/xmlui/bitstream/handle/11356/1364/GOS1.0-words.zip
	cd data && unzip GOS1.0-words.zip

.PHONY: clean
clean:
	rm -v -Rf data/*

.PHONY: build-db
build-db:
	go run main.go app-db

.PHONY: kindle-dict
kindle-dict:
	go run main.go kindle-dict

.PHONY: kindlegen
kindlegen:
	# Works on Kindle previewer installed on OSX:
	"/Applications/Kindle Previewer 3.app/Contents/lib/fc/bin/kindlegen" kindledict/slo-thesaurus.opf -verbose
