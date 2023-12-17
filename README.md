# 페이히어 백엔드 과제
## api개발


페이히어 과제입니다.

## 실행방법
### 
```
> docker-compose up -d payheredb
> createdb.sh
> make build
> docker-compose build payhere
> docker-compose up -d payhere
```

#### m1의 환경인 경우 위의 방법대로 했을 때 실행이 되지 않을 수 있습니다.
```
> vi docker-compose.yml 후 payheredb의 platform 주석을 풀어줍니다.
> docker-compose up -d payheredb
> createdb.sh
> GOOS=Linux make build
> docker-compose build payhere
> docker-compose up -d payhere
```

## API 및 설계
##### 주요 feature와 그 밑에 문장형으로 제 생각 및 의도를 적어두었습니다. 

**Auth**
| Name | URL | Method | Description | Token |
| ------ | ------ | ------ | ------ | ------ |
| register | /auth/register | POST | 회원가입입니다. |-|
| login | /auth/login | POST | 로그인입니다. |-|
| logout | /auh/logout | POST | 로그아웃입니다. |Need|
* register 시 bcrypt를 이용하여 password를 암호화 후 디비에 저장하였습니다. 
* login 시 session 정보를 디비에 저장하고 jwt토큰에 session정보를 함께 저장하였습니다. 토큰을 이용할 때 token의 session정보가 정말 유효한지 디비와 함께 비교합니다. ( 이 과정에서 디비가 아니라 redis나 다른 것을 사용했어도 괜찮았을 것 같습니다.)
* lougout시 디비에 저장된 session정보를 삭제합니다. 

**User**
| Name | URL | Method | Description | Token |
| ------ | ------ | ------ | ------ | ------ |
| NewUser | /user | POST | 사용자가입입니다. |Need|
* 가입을 하여 uid를 발급하는 것과 실제 서비스에 가입하는 것은 좀 다르다고 생각하여 NewUser api를 만들었습니다.
* 카페의 사장님이라고 되어있지만 사장님 이외의 손님이나 다른 타입의 사용자가 나타날 수 있다 생각하였습니다.
* 그 경우에는 계정생성할 때 다루기엔 조금 다른 데이터인 것 같아 NewUser를 통해 사장님으로 가입을 할 수 있게 하였습니다.

**Shop**
| Name | URL | Method | Description | Token |
| ------ | ------ | ------ | ------ | ------ |
| NewShop | /shop | POST | 상점을 생성합니다. |Need|
| GetShopCount | /shop/count | GET | 상점의 개수를 가져옵니다. | Need |
| GetShopList | /shop | GET | 상점의 리스트를 가져옵니다. | Need |
| GetShop | /shop/:shopID | GET | 상점을 가져옵니다. | Need |
| UpdateShop | /shop/:shopID | PUT | 상점을 업데이트합니다. | Need |
| DeleteShop | /shop/:shopID | DELETE | 상점을 삭제합니다. | Need |
* 상점이라는 것은 요구사항에 나와있지 않으나, 상품이라는 것을 생각하였을 때 확장가능성이 있다 생각하여 추가하였습니다.
* user가 단 한개의 상점만 가져 상품들을 가지는 것 보다는 여러 상점을 관리할 수 있고 그 안에 상품을 저장하는 것이 더 확장성있지 않을까 생각했습니다.

**Product**
| Name | URL | Method | Description | Token |
| ------ | ------ | ------ | ------ | ------ |
| NewProduct | /product | POST | 상품을 생성합니다. |Need|
| GetProductCount | /product/count | GET | 상품의 개수를 가져옵니다. | Need |
| GetProductList | /product | GET | 상품의 리스트를 가져옵니다. | Need |
| GetProduct | /product/:productID | GET | 상품을 가져옵니다. | Need |
| UpdateProduct | /product/:productID | PUT | 상품을 업데이트합니다. | Need |
| DeleteProduct | /product/:productID | DELETE | 상품을 삭제합니다. | Need |
* 상품입니다. 상품의 접근 권한은 token안의 uid뿐입니다.
* 리스트의 경우 최소 limit을 10으로 하였습니다. 테이블과 컬럼 정보 기반으로 커서를 생성 후 base64 encoding을 하여 cursor를 생성합니다. 
* 상품 이름으로 검색을 할 때 초성도 가능하여야한다 해서 상품 생성/수정 시 한글 초성을 알려주는 library를 사용하여 데이터베이스에 미리 저장해놓고, 이 정보를 기반으로 search가 되게 하였습니다. (이런 고도의 search는 elasticsearch를 이용하는게 더 좋을까 싶습니다.)
* 카테고리의 경우 고민을 했습니다. 카테고리를 별도의 구조로 빼서 pk를 들고 있게 하는 것도 방법이 될 수 있다고 생각합니다. 하지만 사장님들마다 각자 정의한 카테고리를 사용하고 싶어할 수 있겠다 라는 생각이 들었습니다. 만약 사장님들 마다 고유의 카테고리로 그룹을 해서 리스트를 보여주고 싶다고 한다면 사실 이건 카테고리가 아니라 별도의 productGroup 과 같은 상품의 그룹을 만들어 그 안에 상품을 넣는게 더 좋지 않을까 라는 생각이 들었습니다. 그리고 카테고리는 음료, 음식, 디저트 등 이런 좀 더 큰 범위의 정보가 될 거란 생각이 들었습니다. 만약 이렇게 된다면 당장은 string으로 카테고리를 정리하더라도 나중에 Pk로 변경하는 작업을 하면 된다고 생각하여 일단은 string 필드로 두었습니다. 


**Response**

Response의 경우에는 meta, data 를 키로두고 그 안에 데이터를 넣어야 하는 구조였습니다. 그리고 각 entity마다 data안의 키가 다를 수 있습니다. 그래서 고민하다가 responser라는 인터페이스를 두고 그 인터페이스 함수를 구현한 entity만 reponse에 담을 수 있게 하였습니다. count의 경우에는 단일 int기 때문에 이 경우만 예외로 count라는 키를 만들어서 return하게 했습니다.
이렇게 되면 response로 나가면 안되는 데이터는 controller의 response에서 한번 막을 수 있기 때문에 괜찮은 구조라 생각합니다.


