public class Square
{
    public int X { get; }
    public int Y { get; }
    public int Size { get; }
    public int Right => X + Size;
    public int Bottom => Y + Size;

    public Square(int x, int y, int size)
    {
        X = x;
        Y = y;
        Size = size;
    }
}