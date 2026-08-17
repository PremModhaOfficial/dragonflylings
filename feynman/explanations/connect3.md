the PoolConnetion directly effects how many parallel goroutines can talk to the server without blocking

1. Pooling exists in DBs becouse it is a common but powefull pattern to save time and reuse teh resourtse in this case the connections as making new one eachtime takes extra UNBOUNDED resourses
2. 
