namespace lb1_Backtracking;

public class CLI
{
	private const string Hello = 
	"""
	______               _     _                       _     _               
	| ___ \             | |   | |                     | |   (_)              
	| |_/ /  __ _   ___ | | __| |_  _ __   __ _   ___ | | __ _  _ __    __ _ 
	| ___ \ / _` | / __|| |/ /| __|| '__| / _` | / __|| |/ /| || '_ \  / _` |
	| |_/ /| (_| || (__ |   < | |_ | |   | (_| || (__ |   < | || | | || (_| |
	\____/  \__,_| \___||_|\_\ \__||_|    \__,_| \___||_|\_\|_||_| |_| \__, |
	                                                                    __/ |
	                                                                   |___/ 
	Created by Dubrovin Danila, group: 3388
				
	""";
	
	private const string Help = 
	"""
	Usage: ./lb1_Backtracking [options] [args]

	-h, --help  
		Вывести справку о доступных параметрах и завершить работу.

	-v, --visualize <filename>  
		Выполнить визуализацию решения и записать её в файл filename.  

	-a, --analyze <filename>  
		Запустить анализ производительности (бенчмарк) алгоритма  
		и записать график в файл filename.

	-d, --debug  
		Запустить программу в режиме отладки с выводом подробной информации о ходе выполнения.

	""";
	
	private const string SeeHelp = "See './lb1_Backtracking --help'";

	private readonly Dictionary<string, Action> _flagActions;
	private readonly Dictionary<string, Action<string>> _paramActions;
	private bool _debugMode = false;
	private string? _visualFile;
	private string? _graphFile;

	public CLI()
	{
		_flagActions = new Dictionary<string, Action>
		{
			{"-h",() => Console.WriteLine(Help)},
			{"--help", () => Console.WriteLine(Help)},
			{"-d", () => _debugMode = true},
			{"--debug", () => _debugMode = true}
		};

		_paramActions = new Dictionary<string, Action<string>>
		{
			{"-v", filename => _visualFile = filename},
			{"--visualize", filename => _visualFile = filename},
			{"-a", filename => _graphFile = filename},
			{"--analyze", filename => _graphFile = filename}
		};
	}

	public static void Log(string message, ConsoleColor color = ConsoleColor.White)
	{
		Console.ForegroundColor = color;
		Console.WriteLine(message);
		Console.ResetColor();
	}
	
	private void ParseArguments(string[] args)
	{
		for (int i = 0; i < args.Length; i++)
		{
			if (_flagActions.ContainsKey(args[i]))
			{
				_flagActions[args[i]].Invoke();
				
				if (args[i] == "-h" || args[i] == "--help")
				{
					Environment.Exit(0);
				}
			}
			else if (_paramActions.ContainsKey(args[i]))
			{
				if (i + 1 < args.Length)
				{
					_paramActions[args[i]].Invoke(args[i + 1]);
					i++;
				}
				else
				{
					throw new ArgumentException($"Missed argument for {args[i]}");
				}
			}
			else
			{
				throw new ArgumentException($"Unknown option {args[i]}");
			}
		}
	}
	
	public void Run(string[] args)
	{
		Console.WriteLine(Hello);
		try
		{
			ParseArguments(args);
		}
		catch (Exception ex)
		{
			Log($"Error: {ex.Message}", ConsoleColor.Red);
			Console.WriteLine(SeeHelp);
			return;
		}

		if (_graphFile != null)
		{
			try
			{
				var analyzer = new Analyzer();
				var sizes = new[] { 2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37 };
				var experimentData = analyzer.Benchmark(sizes);
				analyzer.BuildGraph(sizes, experimentData, _graphFile);
			}
			catch (Exception ex)
			{
				Log($"Error while analyzing: {ex.Message}", ConsoleColor.Red);
				return;
			}
		}

		Console.Write("Enter grid size: ");
		if (!int.TryParse(Console.ReadLine(), out int gridSize) || gridSize <= 1)
		{
			Log("Error: Grid size must be a natural number > 1.", ConsoleColor.Red);
			Console.WriteLine(SeeHelp);
			return;
		}

		var partitioner = new SquarePartitioner();
		var result = partitioner.FindOptimalPartition(gridSize, _debugMode);

		Console.WriteLine("---------------------------------------");
		Console.WriteLine($"Minimum squares count: {result.Item1}");
		Console.WriteLine($"Iterations count: {result.Item2}");
		Console.WriteLine("Square positions:");
		foreach (var position in result.Item3)
		{
			Console.WriteLine($"{1 + position[0]} {1 + position[1]} {position[2]}");
		}
		Console.WriteLine("---------------------------------------");

		if (_visualFile != null)
		{
			try
			{
				var visualizer = new Visualizer();
				visualizer.VisualizePartition(result.Item3, gridSize, _visualFile);
			}
			catch (Exception ex)
			{
				Log($"Error while visualizing: {ex.Message}", ConsoleColor.Red);
			}
		}
	}
}
