// gcc 4.7.2 +
// gcc -std=gnu99 -Wall -g -o helloworld_c helloworld_c.c -lpthread

#include <pthread.h>
#include <stdio.h>

int globalNumber=0;
// Note the return type: void*
void* someThreadFunction(){
    printf("Hello from a thread!\n");
    return NULL;
}


void* thread_1(){
	for(int i=0; i<1000000; i++){
		globalNumber++;//I is incremented	
	}
	return NULL;
}
void* thread_2(){
	for (int j=0; j<1000000; j++){
		globalNumber--;
	}
	return NULL;
}


int main(){
    	pthread_t thread1;
	pthread_t thread2;
	pthread_create(&thread1, NULL, thread_1, NULL);
//	pthread_join(thread1,NULL);
   	 // Arguments to a thread would be passed here ---------^
    	pthread_create(&thread2, NULL, thread_2, NULL);

	pthread_join(thread1, NULL);
	pthread_join(thread2, NULL);

	printf("This is the resulting globalNUmber %d \n", globalNumber);
    	return 0;
    
}
