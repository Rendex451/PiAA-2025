using Spectre.Console;

namespace lb1_Backtracking;

public class PerformanceAnalyzer
{
	public int[] Benchmark(int[] sizes)
	{
		var partitioner = new SquarePartitioner();
		var operations = sizes
			.Select(size => partitioner.FindOptimalPartition(size, false).Item2)
			.ToArray();
		return operations;
	}

	public void BuildTable(int[] sizes, int[] operations)
	{
		var table = new Table();
		
		table.AddColumn("Size");
		table.AddColumn("Iterations count");
		
		for (int i = 0; i < sizes.Length; i++)
		{
			var size = sizes[i].ToString();
			var ops = operations[i].ToString();
            
			table.AddRow(size, ops);
		}
		
		AnsiConsole.Write(table);
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
		
		if (!Path.HasExtension(filename))
    	{
        	filename += ".png";
    	}

		var filepath = Path.Combine(Environment.CurrentDirectory, filename);
		plt.Save(filepath, 600, 400);
		CLI.Log($"Graph saved to {filepath}", ConsoleColor.Green);
	}
}