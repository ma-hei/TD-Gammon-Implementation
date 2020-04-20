import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import {HttpParams} from "@angular/common/http";
import { Point } from './point';
import { Checker } from './checker';

@Injectable({
  providedIn: 'root'
})

export class BoardUpdateService {

  private boardBackend = 'http://localhost:1323/GetBoardUpdate';

  constructor(private http: HttpClient) { }
  
  getUpdate(checkersOnPoint: Map<number, Checker[]>, checkersOnBar: Map<number, Checker[]>, playerTurn: number) {
    let params = new HttpParams();
    for (let i=0;i < 24; i++) {
      let checkers = checkersOnPoint.get(i);
      if (checkers !== undefined && checkers.length > 0) {
        params = params.set(i.toString(), checkers.length.toString()+","+checkers[0].player);
      }
    }
    let barCount1 = checkersOnBar.get(0).length;
    let barCount2 = checkersOnBar.get(1).length;
    params = params.set("bar1", barCount1.toString());
    params = params.set("bar2", barCount2.toString());
    params = params.set("playerTurn", playerTurn.toString());
    const options = { params: params };
    return this.http.get(this.boardBackend, options);
  }
}
