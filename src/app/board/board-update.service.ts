import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import {HttpParams} from "@angular/common/http";
import { Point } from './point';
import { Checker } from './checker';

@Injectable({
  providedIn: 'root'
})

export class BoardUpdateService {

  private boardBackend = 'http://localhost:1323';

  constructor(private http: HttpClient) { }
  
  getUpdate(checkersOnPoint: Map<number, Checker[]>) {
    console.log("hi");
    let params = new HttpParams();
    for (let i=0;i < 24; i++) {
      let checkers = checkersOnPoint.get(i);
      if (checkers !== undefined && checkers.length > 0) {
        console.log("adding something");
        console.log(checkers[0].player);
        params = params.set(i.toString(), checkers.length.toString()+","+checkers[0].player);
      }
    }
    const options = { params: params };
    console.log(options);
    return this.http.get(this.boardBackend, options);
  }
}
