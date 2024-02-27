import { UnixDateFormatPipe } from './unix-date-format.pipe';

describe('UnixDateFormatPipe', () => {
  it('create an instance', () => {
    const pipe = new UnixDateFormatPipe();
    expect(pipe).toBeTruthy();
  });
});
