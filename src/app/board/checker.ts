import { BoardDimensions } from './board-dimensions';

export class Checker {
  x: number;
  y: number;
  r: number;
  player: number;
  boardDimensions: BoardDimensions;
  color: string;

  constructor (
    player: number,
    pointId: number,
    nCheckersAlreadyOnPoint: number,
    boardDimensions: BoardDimensions
  ) {
      this.player = player;
      this.x = BoardDimensions.computeCheckerX(pointId, boardDimensions);
      this.y = BoardDimensions.computeCheckerY(pointId, nCheckersAlreadyOnPoint, boardDimensions);
      this.r = BoardDimensions.computeCheckerRadius(boardDimensions);
      this.color = BoardDimensions.getCheckerColor(player);
  }

  moveCheckerToBar(numberCheckersAlreadyOnBar: number, boardDimensions: BoardDimensions) {
      this.x = BoardDimensions.computeCheckerOnBarX(boardDimensions);
      this.y = BoardDimensions.computeCheckerOnBarY(this.player==1, numberCheckersAlreadyOnBar, boardDimensions); 
  }  

  moveCheckerToPoint(nCheckersAlreadyOnPoint: number, pointId: number, boardDimensions: BoardDimensions) {
     this.x = BoardDimensions.computeCheckerX(pointId, boardDimensions);
     this.y = BoardDimensions.computeCheckerY(pointId, nCheckersAlreadyOnPoint, boardDimensions);
  }

  bearCheckerOff(boardDimensions: BoardDimensions) {
     this.x = BoardDimensions.computeCheckerBearOffX(boardDimensions);
  }
}
