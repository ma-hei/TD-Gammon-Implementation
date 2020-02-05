import { Component, OnInit, Input } from '@angular/core';
import { Point } from './point';
import { Checker } from './checker';
import { BoardDimensions } from './board-dimensions';
import { BoardUpdateService } from './board-update.service';

@Component({
  selector: 'app-board',
  templateUrl: './board.component.html',
  styleUrls: ['./board.component.css']
})

export class BoardComponent implements OnInit {

  @Input() boardWidth: number;
  @Input() boardHeight: number;
  leftOffset: number;
  pointHeight: number;
  pointSpacing: number;
  pointWidth: number;
  points: Point[] = new Array();
  checkers: Checker[] = new Array();
  boardDimensions: BoardDimensions;
  checkersOnPoint: Map<number, Checker[]> = new Map<number, Checker[]>();

  constructor( private boardUpdateService: BoardUpdateService) {
    this.leftOffset = 10;
    this.pointHeight = 200;
  }

  ngOnInit() {
    this.boardDimensions = new BoardDimensions(this.boardWidth, this.boardHeight);
    for (var i = 0; i < 24; i++){
      this.points[i] = new Point(i, this.boardDimensions); 
      this.checkersOnPoint.set(i, new Array());
    } 
    this.putCheckersOnPoint(5, 0, 1);
    this.putCheckersOnPoint(3, 4, 2);
    this.putCheckersOnPoint(5, 6, 2);
    this.putCheckersOnPoint(2, 11, 1);
    this.putCheckersOnPoint(5, 12, 2);
    this.putCheckersOnPoint(3, 16, 1);
    this.putCheckersOnPoint(5, 18, 1);
    this.putCheckersOnPoint(2, 23, 2);
    this.boardUpdateService.getUpdate(this.checkersOnPoint).subscribe((data) => console.log(data)); 
  }

  putCheckersOnPoint(nCheckers: number, pointId: number, player: number) {
    let currentCheckersOnPoint = this.checkersOnPoint.get(pointId);
    for (let i = 0; i < nCheckers; i++) {
      let checker = new Checker(player, pointId, currentCheckersOnPoint.length, this.boardDimensions);
      this.checkers.push(checker);
      this.checkersOnPoint.get(pointId).push(checker);
    }
  }

  mouseEnter(id: number) {
  }
}
