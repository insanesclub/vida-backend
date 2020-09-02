#### vida-backend

Back-end module of VIDA(Visualization tool for Instagram Data Analysis)



0. Index

   - /details
   
   

3. /details

   음식점들의 상세 정보(상호명, 전화번호, 도로명 주소, 평균 가격)를 검색한 결과를 반환합니다.

   | method |     path     |     request     |               response               |
   | :----: | :----------: | :-------------: | :----------------------------------: |
   | `GET`  | /api/details | (string) 검색어 | (Array&lt;JSON&gt;) 음식점 상세 정보 |
   
   - Query param 예시

     `/api/details?tag=finedining`

   - Response body 예시

     ![sample](https://user-images.githubusercontent.com/29545214/92040089-268d2500-edb1-11ea-9ff3-0b01a135ce29.png)

     - place_name: (string) 상호명
     - road_address_name: (string) 도로명 주소
     - phone: (string) 전화번호
     - average_price: (string) 평균 가격 (내림차순으로 정렬되며, 가격 정보가 없을 경우 빈 문자열입니다.)
