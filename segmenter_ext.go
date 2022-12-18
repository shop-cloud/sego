/*
 * @Author: lwnmengjing<lwnmengjing@qq.com>
 * @Date: 2022/12/18 22:25:57
 * @Last Modified by: lwnmengjing<lwnmengjing@qq.com>
 * @Last Modified time: 2022/12/18 22:25:57
 */

package sego

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"strconv"
)

func (seg *Segmenter) LoadDictionaryFromReader(readers ...io.Reader) {
	seg.dict = NewDictionary()
	for i := range readers {
		log.Printf("载入sego词典 %d", i)
		reader := bufio.NewReader(readers[i])
		var text string
		var freqText string
		var frequency int
		var pos string

		// 逐行读入分词
		for {
			size, _ := fmt.Fscanln(reader, &text, &freqText, &pos)

			if size == 0 {
				// 文件结束
				break
			} else if size < 2 {
				// 无效行
				continue
			} else if size == 2 {
				// 没有词性标注时设为空字符串
				pos = ""
			}

			// 解析词频
			var err error
			frequency, err = strconv.Atoi(freqText)
			if err != nil {
				continue
			}

			// 过滤频率太小的词
			if frequency < minTokenFrequency {
				continue
			}

			// 将分词添加到字典中
			words := splitTextToWords([]byte(text))
			token := Token{text: words, frequency: frequency, pos: pos}
			seg.dict.addToken(token)
		}
	}

	// 计算每个分词的路径值，路径值含义见Token结构体的注释
	logTotalFrequency := float32(math.Log2(float64(seg.dict.totalFrequency)))
	for i := range seg.dict.tokens {
		token := &seg.dict.tokens[i]
		token.distance = logTotalFrequency - float32(math.Log2(float64(token.frequency)))
	}

	// 对每个分词进行细致划分，用于搜索引擎模式，该模式用法见Token结构体的注释。
	for i := range seg.dict.tokens {
		token := &seg.dict.tokens[i]
		segments := seg.segmentWords(token.text, true)

		// 计算需要添加的子分词数目
		numTokensToAdd := 0
		for iToken := 0; iToken < len(segments); iToken++ {
			if len(segments[iToken].token.text) > 0 {
				numTokensToAdd++
			}
		}
		token.segments = make([]*Segment, numTokensToAdd)

		// 添加子分词
		iSegmentsToAdd := 0
		for iToken := 0; iToken < len(segments); iToken++ {
			if len(segments[iToken].token.text) > 0 {
				token.segments[iSegmentsToAdd] = &segments[iToken]
				iSegmentsToAdd++
			}
		}
	}

	log.Println("sego词典载入完毕")
}
