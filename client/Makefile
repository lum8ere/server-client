.PHONY: build

run-old:
	py client.py

run:
	python -m client.main

build:
	pyinstaller --onefile -w --uac-admin client.py 

clean:
	rm -r ./dist
	rm -r ./build
	rm client.spec