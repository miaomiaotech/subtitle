# pip install git+https://github.com/weaming/thunder-subtitle-2024.git
function download-subitltes -a ext
    for x in (ls *.$ext)
        echo $x
        thunder-subtitle -i 1 $x
    end
    echo 1. open https://srttran.leavesc.com to translate
    echo 2. copy results to the same directory
    echo 3. run merge-all-subtitles
end

# go install -v -ldflags "-s -w" github.com/miaomiaotech/subtitle/cmd/merge-subtitles@latest
function merge-all-subtitles
    CHINESE_FIRST=1 merge-subtitles (ls *.srt | sort -n)

    echo '---------------'
    echo clean srt files...

    mkdir archive
    mv *.zh-CN.srt ./archive/

    for x in (ls *.merged.srt)
        set -l origin (echo $x | sed 's/\.zh-CN\.merged//g')
        echo $origin
        mv $origin ./archive/
        mv $x $origin
    end
end
