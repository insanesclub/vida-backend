#### vida-backend

Back-end module of VIDA(Visualization tool for Instagram Data Analysis)



0. Index

   - Fake user detector
- Influencer intervention measurement
   
   


1. Fake user detector (고의적 도배 이용자 감별기)

   같은 주제(해시태그)에 대해 **고의적으로 같은 게시물을 올리는 사용자**를 찾아내고, 그 결과를 시각화하는 서비스입니다.

   매일매일 해당 해시태그를 포함한 게시물을 업로드하는 사용자들의 빈도를 측정합니다.

   만약 오늘 `#신촌맛집`으로 `user_A`라는 사용자가 3개의 게시물을, `user_B`라는 사용자가 2개의 게시물을 업로드했다면 다음과 같은 정보가 저장됩니다.

   ![sample](https://user-images.githubusercontent.com/29545214/88458101-69accc00-cec6-11ea-833b-f05fed461b05.png)

   

   대시보드에서는 두가지 그래프로 데이터를 시각화합니다.

   

   첫째는 오늘의 결과 그래프입니다.

   x축은 사용자의 닉네임, y축은 올린 게시물의 수가 될 것입니다.

   오늘 측정한 결과를 내림차순 정렬하고 **막대 그래프**로 보여줍니다.

   백엔드에서는 `tag`와 `date` 값을 받아 `uploads` 필드의 값을 반환하는 API를 제공합니다.

   

   둘째는 누적 결과 그래프입니다.

   DB에 저장된 모든 날의 누적 결과를 모아, 지금까지 해당 해시태그로 꾸준히 게시물을 올린 사용자들의 업로드 횟수를 보여줍니다.

   보여주는 방식은 첫째와 같습니다.

   백엔드에서는 `tag` 값을 받아 `uploads` 필드들의 누적 합산 결과를 반환하는 API를 제공합니다.

   

2. Influencer intervention measurement (인플루언서 개입 정도 측정기)

   위치 태그 검색을 통해 해당 상호에 **인플루언서들이 개입한 정도**를 시각화하는 서비스입니다.

   이를 통해 초반에 게시물을 올린 사용자들의 평균 팔로워 수가 지나치게 높은, 인플루언서와의 협찬 계약을 통해 유명세를 탄 상호를 가려낼 수 있을 것으로 보입니다.

   만약 오늘 `#돈부리파스타` 를 위치태그로 한 게시물을 올린 사용자 `user_A` 의 팔로워 수가 25만명, `user_B` 의 팔로워 수가 19만명이었다면 다음과 같은 정보가 저장됩니다.

   ![sample](https://user-images.githubusercontent.com/29545214/88458152-ea6bc800-cec6-11ea-800e-a22f0d2d353a.png)

   

   `// TODO: 2번 아이디어에 대한 시각화 방법 구체화`

