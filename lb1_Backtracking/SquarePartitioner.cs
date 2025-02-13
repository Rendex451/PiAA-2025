public class SquarePartitioner
{
    private bool _debugMode = false;
    public Tuple<int, int, List<int[]>> FindOptimalPartition(int gridSize, bool debugMode)
    {
        _debugMode = debugMode;
        
        int operationCounter = 0;
        int squareSize;
        int newGridSize = ScaleSquareSize(gridSize, out squareSize);
        int bestCount = 2 * newGridSize + 1;
        List<Square> squares = PlaceInitialSquares(newGridSize);
        List<Square> bestSolution = new List<Square>();
        int initialOccupiedArea = squares[0].Size * squares[0].Size + 2 * squares[1].Size * squares[1].Size;
        int startX = squares[2].Bottom, startY = squares[2].X;
        
        Backtrack(
            squares, 
            bestSolution, 
            initialOccupiedArea, 
            3, 
            startX, startY, 
            newGridSize, 
            ref bestCount, 
            ref operationCounter);

        List<int[]> formattedSolution = bestSolution
            .Select(square => 
                new int[] { square.X * squareSize, square.Y * squareSize, square.Size * squareSize }
            )
            .ToList();
        
        _debugMode = false;
        return Tuple.Create(bestCount, operationCounter, formattedSolution);
    }
    
    private bool IsPositionOccupied(List<Square> squares, int x, int y)
    {
        return squares.Any(square =>
            x >= square.X && x < square.Right && y >= square.Y && y < square.Bottom);
    }

    private void Backtrack(
        List<Square> squares, 
        List<Square> bestSolution,
        int occupiedArea,
        int currentCount,
        int startX, int startY,
        int gridSize,
        ref int bestCount,
        ref int operationCounter
    ) 
    {
        if (_debugMode)
            Console.WriteLine("[New Backtrack Iteration]");
        
        operationCounter++;
        if (occupiedArea == gridSize * gridSize)
        {
            UpdateBestSolution(squares, bestSolution, currentCount, ref bestCount);
            return;
        }

        for (int x = startX; x < gridSize; x++)
        {
            for (int y = startY; y < gridSize; y++)
            {
                if (IsPositionOccupied(squares, x, y))
                    continue;

                int maxSize = CalculateMaxSquareSize(squares, x, y, gridSize);
                if (maxSize <= 0)
                    continue;

                TryPlacingSquares(squares,
                    bestSolution,
                    occupiedArea,
                    currentCount,
                    x, y,
                    gridSize,
                    maxSize,
                    ref bestCount,
                    ref operationCounter);
                
                return;
            }
            startY = 0;
        }
    }
    
    private void UpdateBestSolution(List<Square> squares, List<Square> bestSolution, int currentCount, ref int bestCount)
    {
        if (currentCount < bestCount)
        {
            bestCount = currentCount;
            bestSolution.Clear();
            bestSolution.AddRange(squares);
        }
    }
    
    private int CalculateMaxSquareSize(List<Square> squares, int x, int y, int gridSize)
    {
        int maxSize = Math.Min(gridSize - x, gridSize - y);
    
        foreach (var square in squares)
        {
            if (square.Right > x && square.Y > y)
                maxSize = Math.Min(maxSize, square.Y - y);
            else if (square.Bottom > y && square.X > x)
                maxSize = Math.Min(maxSize, square.X - x);
        }
    
        return maxSize;
    }
    
    private void TryPlacingSquares(
        List<Square> squares, 
        List<Square> bestSolution,
        int occupiedArea,
        int currentCount,
        int x, int y, 
        int gridSize,
        int maxSize,
        ref int bestCount,
        ref int operationCounter)
    {
        for (int size = maxSize; size >= 1; size--)
        {
            var newSquare = new Square(x, y, size);
            var newOccupiedArea = occupiedArea + size * size;
        
            squares.Add(newSquare);
            if (_debugMode)
                Console.WriteLine($"[New Square ({newSquare.X} {newSquare.Y} {newSquare.Size})]");
            if (newOccupiedArea == gridSize * gridSize)
            {
                UpdateBestSolution(squares, bestSolution, currentCount + 1, ref bestCount);
            }
            else if (currentCount + 1 < bestCount)
            {
                Backtrack(squares,
                    bestSolution,
                    newOccupiedArea,
                    currentCount + 1,
                    x, y, gridSize,
                    ref bestCount,
                    ref operationCounter);
            }
        
            squares.RemoveAt(squares.Count - 1);
        }
    }

    private List<Square> PlaceInitialSquares(int gridSize)
    {
        int mainSquareSize = (gridSize + 1) / 2;
        int subSquaresSize = gridSize / 2;
        
        if (_debugMode)
        {
            Console.WriteLine($"[New Square (0 0 {mainSquareSize})]");
            Console.WriteLine($"[New Square (0 {mainSquareSize} {subSquaresSize})]");
            Console.WriteLine($"[New Square ({mainSquareSize} 0 {subSquaresSize})]");
        }
        
        return new List<Square>
        {
            new Square(0, 0, mainSquareSize),
            new Square(0, mainSquareSize, subSquaresSize),
            new Square(mainSquareSize, 0, subSquaresSize)
        };
    }

    private int ScaleSquareSize(int gridSize, out int squareSize)
    {
        int maxDivisor = 1;
        for (int i = gridSize / 2; i >= 1; i--)
        {
            if (gridSize % i == 0)
            {
                maxDivisor = i;
                break;
            }
        }
        squareSize = maxDivisor;
        return gridSize / maxDivisor;
    }
}