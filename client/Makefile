.PHONY: build

run:
	py client.py

build:
	pyinstaller --onefile -w --uac-admin client.py 

clean:
	rm -r ./dist
	rm -r ./build
	rm client.spec