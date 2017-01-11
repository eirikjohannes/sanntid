# Python 3.3.3 and 2.7.6
# python helloworld_python.py

from threading import Thread

i = 0

def someThreadFunction():
    print("Hello from a thread!")

# Potentially useful thing:
#   In Python you "import" a global variable, instead of "export"ing it when you declare it
#   (This is probably an effort to make you feel bad about typing the word "global")
    global i

def thread_1():
	global i
	for x in range(0,1000000):
		i+=1 

def thread_2():
	global i
	for y in range(0,1000000):
		i-=1

def main():

	someThread = Thread(target = thread_1, args = (),)
	someThread2 = Thread(target=thread_2,args=(),)
	someThread.start()
	someThread2.start()
        someThread.join()
	someThread2.join()
    	print "This iz our numbah:",i


main()
