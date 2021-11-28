import http from 'k6/http';

export const options = {
  vus: 10,
  duration: '1h',
}

const randomString = () => {
    var result           = '';
    var characters       = 'abcdefghijklmnopqrstuvwxyzกขฃคฅฆงจฉชซฌญฎฏฐฑฒณดตถทธนบปผฝพฟภมยรฤฤๅลฦฦๅวศษสหฬอฮ';
    for ( var i = 0; i < 5; i++ ) {
      result += characters.charAt(Math.floor(Math.random() * 
        characters.length));
   }
   return result;
    
}

export default function() {
    http.get('http://host.docker.internal:8008/books/get?name='+randomString())
}