下载模型:
```
mkdir model
cd model
wget https://huggingface.co/TheBloke/openchat_3.5-GGUF/blob/main/openchat_3.5.Q8_0.gguf
wget https://huggingface.co/openchat/openchat_3.5/blob/main/tokenizer.json
mv tokenizer.json openchat_3.5_tokenizer.json
cd ..
```
运行:
```
bash run.sh
```
