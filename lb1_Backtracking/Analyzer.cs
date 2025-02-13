namespace lb1_Backtracking;

public class Analyzer
{
    public int[] Benchmark(int[] sizes)
    {
        var partitioner = new SquarePartitioner();
        var operations = sizes
            .Select(size => partitioner.FindOptimalPartition(size, false).Item2)
            .ToArray();
        return operations;
    }
    
    public void BuildGraph(int[] sizes, int[] operations, string filename)
    {
        if (sizes.Length != operations.Length)
        {
            throw new ArgumentException("The number of sizes must be equal to the number of operations.");
        }
        
        var plt = new ScottPlot.Plot();

        plt.Add.Scatter(sizes, operations);
        plt.Title("Growth of Iterations Count on Prime Sizes");
        plt.XLabel("Grid Size");
        plt.YLabel("Iterations Count");

        plt.Save($"{filename}.png", 600, 400);
    }
}