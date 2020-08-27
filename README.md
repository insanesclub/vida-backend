#### vida-backend

Back-end module of VIDA(Visualization tool for Instagram Data Analysis)



- Server domain:Port number

  Public DNS(IPv4): `http://ec2-3-35-24-235.ap-northeast-2.compute.amazonaws.com:8080`

  

0. Index

   - Fake user detector
   - Influencer intervention measurement
   - prices
   
   


1. Fake user detector (고의적 도배 이용자 감별기)

   같은 주제(해시태그)에 대해 **고의적으로 같은 게시물을 올리는 사용자**를 찾아내고, 그 결과를 시각화하는 서비스입니다.

   매일매일 해당 해시태그를 포함한 게시물을 업로드하는 사용자들의 빈도를 측정합니다.

   | method |     path     |      request      |                         response                          |
   | :----: | :----------: | :---------------: | :-------------------------------------------------------: |
   | `GET`  | /api/uploads | (string) 해시태그 | (JSON) 금일 해당 해시태그로 게시물을 올린 사용자들의 빈도 |

   - Query string 예시
   
     `tag=신촌맛집`

   - Response body 예시

     만약 오늘 `#신촌맛집`으로 `user_A`라는 사용자가 3개의 게시물을, `user_B`라는 사용자가 2개의 게시물을 업로드했다면 다음과 같은 내용의 JSON이 반환됩니다.
   
     ![sample](https://user-images.githubusercontent.com/29545214/88988534-21513c00-d314-11ea-87d8-ecee6c18c2e7.png)

   

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



3. prices

   음식점들의 상세 정보(상호명, 전화번호, 도로명 주소, 메뉴와 가격 리스트)를 알아내는 서비스입니다.

   이를 통해 특수 작물을 판매하는 음식점을 파악하기 수월해질 것으로 보입니다.

   | method |    path     |    request    |        response         |
   | :----: | :---------: | :-----------: | :---------------------: |
   | `GET`  | /api/prices | (string) 날짜 | (JSON) 음식점 상세 정보 |

   - Query param 예시

     `/api/prices/date=20200828`

   - Response body 예시

     ![sample](https://user-images.githubusercontent.com/29545214/91490868-1b884f80-e8ee-11ea-870f-28ad467c5b60.png)

     

     - date: (timestamp) 날짜
     - place_info: (Array) 해당 일자에 수집한 음식점 정보
     - place_name: (string) 상호명
     - phone: (string) 전화번호
     - road_address-name: (string) 도로명 주소
     - menu_list: (Array) 메뉴와 가격 정보
     - menu_name: (string) 메뉴 이름
     - menu_price: (number) 메뉴 가격