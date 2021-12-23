# Request Counter



  

  

<br  />

<p  align="center">

  

<a  href="https://github.com/Yangiboev/request-counter"></a>

  

<h3  align="center">Request Counter</h3>

<p  align="center">

<br  />

<a  href="https://github.com/Yangiboev/request-counter"><strong>Explore the docs »</strong></a>

  

<br  />

<!-- TABLE OF CONTENTS -->

<details  open="open">

  

<summary>Table of Contents</summary>

<ol>

<li><a  href="#about-the-repo">About The Repo</a></li>

<li><a  href="#getting-started">Getting Started</a><ul>

<li><a  href="#installation">Installation</a></li>

<li><a  href="#contact">Contact</a></li>

</ol>

  

</details>

  

  

<!-- ABOUT THE PROJECT -->

  

## About The Repo

  

This is a simple server which is being develeped by Dilmurod
  
* Our time should be focused on creating something amazing. A project that solves a problem and helps others

* I should element DRY principles to the rest of your life :smile:

* I assume that  the web server is synchronous, if the other case (asynchronous), then I would have to guard my counter using a mutex or atomic in order to prevent my server from being hit with race-condition bugs. However, I believe that locking would be expensive if we get hit a high-performance scenario.

* After considering a high-performance scenarios that thousands of requests could come in at the very same time, I thought that I need proper synchronization. The concurrency model of golang makes it easy to serialize all accesses to the application's state: no (explicit) mutexes are required.






<!-- GETTING STARTED -->

  

## Getting Started

  
  

This is an example of how you may give instructions on setting up your project locally.

  

To get a local copy up and running follow these simple example steps.



### Installation


1. Clone the repo

  

```sh
git clone https://github.com/Yangiboev/request-counter.git
```

If you run code locally please, make sure that you have golang version 1.15 or above.

2. In order to run locally

  

```sh
go run cmd/main.go
```



1. Fork the Project

  

2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)

  

3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)

  

4. Push to the Branch (`git push origin feature/AmazingFeature`)

  

5. Open a Pull Request

  
  

<!-- CONTACT -->

  

## Contact

  

  

Dilmurod Yangiboev - [@icon_me](dilmurod.yangiboev@gmail.com) - dilmurod.yangiboev@gmail.com

  

Repo Link: [https://github.com/Yangiboev/request-counter](https://github.com/Yangiboev/request-counter)

  
