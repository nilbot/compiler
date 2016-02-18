#!$(which ruby)
#This script generate $word_count random words and save them into file
#"words.test"
def coin_flip(bias, head, tail)
	outcome = Random.new.rand
	return head if outcome < bias
	return tail if outcome >= bias
end
def format(rst)
	str = ""
	rst.each {
		|w|
		str+=w
		str+="\n"
	}
	str
end

allowed = (35..126).to_a # ascii alphabet excluding space, ! and " (32,33,34)
alphabets = ((65..90).to_a << (97..122).to_a).flatten! # A-Za-z
length_range = (1..64).to_a # length of a word
word_count = 32000 # count of all words
biased = 0.75 # probability threshold for guaranteed positive testcase


rst = []
word_count.times do
	# pool = coin_flip(biased,alphabets, allowed)
	pool = alphabets
	len = Random.new.rand(length_range.count)
	word = ""
	0.upto(len) do
		word+=(pool[Random.new.rand(pool.count)]).chr
	end
	rst.push(word)
end

f = File.open("words.test", "w")
f.puts format(rst)
f.close