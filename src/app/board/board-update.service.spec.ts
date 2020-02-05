import { TestBed, inject } from '@angular/core/testing';

import { BoardUpdateService } from './board-update.service';

describe('BoardUpdateService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [BoardUpdateService]
    });
  });

  it('should be created', inject([BoardUpdateService], (service: BoardUpdateService) => {
    expect(service).toBeTruthy();
  }));
});
