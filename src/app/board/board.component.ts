import { Component, OnInit, Input } from '@angular/core';
import { Subject } from 'rxjs';
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
  checkersOnBar: Map<number, Checker[]> = new Map<number, Checker[]>();
  subject: Subject<number>;
  nTurns: number;

  constructor( private boardUpdateService: BoardUpdateService) {
    this.leftOffset = 10;
    this.pointHeight = 200;
  }

  ngOnInit() {
    this.checkersOnBar.set(0, new Array());
    this.checkersOnBar.set(1, new Array());
    this.boardDimensions = new BoardDimensions(this.boardWidth, this.boardHeight);
    for (var i = 0; i < 24; i++){
      this.points[i] = new Point(i, this.boardDimensions); 
      this.checkersOnPoint.set(i, new Array());
    } 
    this.putNewCheckersOnPoint(5, 0, 1);
    this.putNewCheckersOnPoint(3, 4, 2);
    this.putNewCheckersOnPoint(5, 6, 2);
    this.putNewCheckersOnPoint(5, 12, 2);
    this.putNewCheckersOnPoint(2, 11, 1);
    this.putNewCheckersOnPoint(3, 16, 1);
    this.putNewCheckersOnPoint(5, 18, 1);
    this.putNewCheckersOnPoint(2, 23, 2);
    //this.putCheckerOnBar(4);
    //this.putCheckerOnBar(4);
    //this.putCheckerOnBar(22);
    //this.putCheckerOnBar(22);
    this.nTurns = 0;
    //   this.updateBoard(data["returnState"])
    this.subject = new Subject<number>();
    this.subject.subscribe(value => {
        console.log("done with turn " + value);
        if (value < 5) {
            setTimeout(this.makeTurn(value), 5000);
        }
    });
    this.makeTurn(0); 
  }

  makeTurn(turnNumber: number) {
      let playerTurn = (turnNumber % 2) + 1;
      console.log("getting random dice roll and board update for player " + playerTurn);
      this.boardUpdateService.getUpdate(this.checkersOnPoint, this.checkersOnBar, playerTurn).subscribe((data) => {
         console.log(data["returnState"]);
         this.updateBoard(data["returnState"], playerTurn);
         this.subject.next(turnNumber + 1);
      });
  }

  performMove(move: string, player: number) {
      console.log("performing move: " + move);
      let idlm = move.indexOf("lm");
      let m = move.substring(idlm+4);
      let idxDot = m.indexOf(".")
      let first = m.substring(0, idxDot)
      let second = m.substring(idxDot+1)
      let idxh = second.indexOf("h");
      let toPoint = -1;
      let isHit = false;
      if (idxh == -1) {
          toPoint = parseInt(second);
      } else {
          toPoint = parseInt(second.substring(0, idxh));
          isHit = true;
      }
      if (first == "b") {
          this.moveCheckerFromBar(toPoint, player);
      } else {
          let fromPoint = parseInt(first);
          if (toPoint < 24) {
              if (isHit) {
                  this.putCheckerOnBar(toPoint);
              }
              this.moveCheckerFromPointToPoint(fromPoint, toPoint);
          } else {
              this.bearCheckerOff(fromPoint);          
          }
      }
  }

  updateBoard(newState: string, playerTurn: number) {
     let idx1 = newState.indexOf("lm1");
     let idx2 = newState.indexOf("lm2");
     let idx3 = newState.indexOf("lm3");
     let idx4 = newState.indexOf("lm4");
     let move1 = newState.substring(idx1, idx2);
     let move2 = newState.substring(idx2, idx3);
     let move3 = newState.substring(idx3, idx4);
     let move4 = newState.substring(idx4);
     if (move1 !== "lm1:") {
         this.performMove(move1, playerTurn);
     }
     if (move2 !== "lm2:") {
         this.performMove(move2, playerTurn);
     }
     if (move3 !== "lm3:") {
         this.performMove(move3, playerTurn);
     }
     if (move4 !== "lm4:") {
         this.performMove(move4, playerTurn);
     }
  }

  bearCheckerOff(fromPoint) {
      let nCheckersOnPointOfPlayer = this.checkersOnPoint.get(fromPoint).length;
      let checker = this.checkersOnPoint.get(fromPoint)[nCheckersOnPointOfPlayer -1];
      this.checkersOnPoint.get(fromPoint).splice(-1,1);
  }

  moveCheckerFromPointToPoint(fromPoint: number, toPoint: number) {
      let nCheckersOnPointOfPlayer = this.checkersOnPoint.get(fromPoint).length;
      let checker = this.checkersOnPoint.get(fromPoint)[nCheckersOnPointOfPlayer -1];
      this.checkersOnPoint.get(fromPoint).splice(-1,1);
      this.putCheckerOnPoint(checker, toPoint);
  }

  moveCheckerFromBar(toPoint: number, player: number) {
      let nCheckersOnBarOfPlayer = this.checkersOnBar.get(player - 1).length;
      let checker = this.checkersOnBar.get(player-1)[nCheckersOnBarOfPlayer - 1];
      this.putCheckerOnPoint(checker, toPoint);
      this.checkersOnBar.get(player - 1).splice(-1,1);
  }

  putCheckerOnPoint(checker: Checker, pointId: number) {
     let currentCount = this.checkersOnPoint.get(pointId).length;
     checker.moveCheckerToPoint(currentCount, pointId, this.boardDimensions);
     this.checkersOnPoint.get(pointId).push(checker);
  }

  putCheckerOnBar(fromPoint: number) {
     let nCheckersOnPoint = this.checkersOnPoint.get(fromPoint).length
     let checker = this.checkersOnPoint.get(fromPoint)[nCheckersOnPoint-1];
     let nCheckersOnBarOfPlayer = this.checkersOnBar.get(checker.player - 1).length;
     checker.moveCheckerToBar(nCheckersOnBarOfPlayer, this.boardDimensions);
     this.checkersOnBar.get(checker.player - 1).push(checker);
     this.checkersOnPoint.get(fromPoint).splice(-1,1);
  }

  putNewCheckersOnPoint(nCheckers: number, pointId: number, player: number) {
    let currentCheckersOnPoint = this.checkersOnPoint.get(pointId);
    for (let i = 0; i < nCheckers; i++) {
      let checker = new Checker(player, pointId, currentCheckersOnPoint.length, this.boardDimensions);
      this.checkers.push(checker);
      this.checkersOnPoint.get(pointId).push(checker);
    }
  }

  mouseEnter(id: number) {
    console.log(id);
  }
}
