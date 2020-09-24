#### vida-backend

Back-end module of VIDA(Visualization tool for Instagram Data Analysis)



0. Index

   - /details
   
   

1. /details

   음식점들의 상세 정보(상호명, 전화번호, 도로명 주소, 평균 가격)를 검색한 결과를 반환합니다.

   | method |     path     |     request     |               response               |
   | :----: | :----------: | :-------------: | :----------------------------------: |
   | `GET`  | /api/details | (string) 검색어 | (Array&lt;JSON&gt;) 음식점 상세 정보 |

   - Query param 예시

     `/api/details?tag=서울맛집추천`

   - Response body 예시

     ```json
     [
       {
         "place_name":"오늘",
         "road_address_name":"서울 용산구 장문로 60",
         "phone":"02-792-1054",
         "average_price":"108600"
       },
       {
         "place_name":"델비노",
         "road_address_name":"서울 광진구 워커힐로 177",
         "phone":"02-2022-0111",
         "average_price":"85000"
       },
       {
         "place_name":"콩두",
         "road_address_name":"서울 중구 덕수궁길 116-1",
         "phone":"02-722-7002",
         "average_price":"45050"
       },
       {
         "place_name":"조인바이트",
         "road_address_name":"서울 강남구 선릉로112길 13",
         "phone":"02-3444-3611",
         "average_price":"38000"
       },
       {
         "place_name":"제주돈사돈",
         "road_address_name":"부산 수영구 수영로 765-1",
         "phone":"051-756-3001",
         "average_price":"35333"
       },
       {
         "place_name":"150",
         "road_address_name":"서울 강남구 도산대로46길 8",
         "phone":"",
         "average_price":"32600"
       },
       {
         "place_name":"꿈꾸는포장마차",
         "road_address_name":"서울 동작구 동작대로7길 35",
         "phone":"02-523-1020",
         "average_price":"32500"
       },
       {
         "place_name":"북한산우동집",
         "road_address_name":"경기 고양시 덕양구 북한산로 639",
         "phone":"02-354-0818",
         "average_price":"32280"
       }
       ...
     ]
     ```

     - place_name: (string) 상호명
     - road_address_name: (string) 도로명 주소
     - phone: (string) 전화번호
     - average_price: (string) 평균 가격 (내림차순으로 정렬되며, 가격 정보가 없을 경우 빈 문자열입니다.)
