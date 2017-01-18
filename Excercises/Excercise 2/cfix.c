// gcc 4.7.2 +
// gcc -std=gnu99 -Wall -g -o helloworld_c helloworld_c.c -lpthread

#include <pthread.h>
#include <stdio.h>

pthread_mutex_t globalNumberLock;
int globalNumber=0;
// Note the return type: void*
void* thread_1(){
	for(int i=0; i<1000000; i++){
		pthread_mutex_lock(&globalNumberLock);
		globalNumber++;//I is incremented
		pthread_mutex_unlock(&globalNumberLock);	
	}
	return NULL;
}
void* thread_2(){
	for (int j=0; j<1000001; j++){
		pthread_mutex_lock(&globalNumberLock);
		globalNumber--;
		pthread_mutex_unlock(&globalNumberLock);
	}
	return NULL;
}


int main(){
	pthread_mutex_init(&globalNumberLock,NULL);

    	pthread_t thread1;
	pthread_t thread2;
	pthread_create(&thread1, NULL, thread_1, NULL);
//	pthread_join(thread1,NULL);
   	 // Arguments to a thread would be passed here ---------^
    	pthread_create(&thread2, NULL, thread_2, NULL);

	pthread_join(thread1, NULL);
	pthread_join(thread2, NULL);
	pthread_mutex_destroy(&globalNumberLock);
	printf("This is the resulting globalNUmber %d \n", globalNumber);
    	return 0;
    
}
