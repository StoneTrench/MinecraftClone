# mv ./_out/logs ./logs_backup
# rm -rf ./_out/*
# mv ./logs_backup ./_out/logs

# mkdir ./_out/assets/
# mkdir ./_out/assets/shaders/
# slangc assets/shaders/shader.slang -target spirv -profile spirv_1_3 -emit-spirv-directly -fvk-use-entrypoint-name -entry vertMain -entry fragMain -o _out/assets/shaders/shader.spv
for file in assets/shaders/*.slang; do
    filename=$(basename "$file" .slang)
    
    slangc "$file" \
        -target spirv \
        -profile spirv_1_3 \
        -emit-spirv-directly \
        -fvk-use-entrypoint-name \
        -entry vertMain \
        -entry fragMain \
        -o "_out/assets/shaders/${filename}.spv"
done

cp -rf assets/textures/ _out/assets/textures/

# for file in include/*; do
#     filename=$(basename "$file")
# 		cp $file ./_out/$filename
# done