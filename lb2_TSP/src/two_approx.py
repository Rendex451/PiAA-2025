import sys
from argparse import ArgumentParser
from heapq import heappush, heappop


class DebugLogger:
    def __init__(self, debug=False):
        self.debug = debug

    def log(self, message):
        if self.debug:
            print(f"[DEBUG] {message}", file=sys.stderr)

    def log_mst(self, mst):
        if self.debug:
            print("\nМинимальное остовное дерево (MST):", file=sys.stderr)
            for i, edges in enumerate(mst):
                print(f"Вершина {i} соединена с:", end=" ", file=sys.stderr)
                for v, w in edges:
                    print(f"{v}(вес {w})", end=" ", file=sys.stderr)
                print(file=sys.stderr)

    def log_path(self, path, graph):
        if self.debug:
            print("\nПостроенный путь:", " -> ".join(map(str, path)), file=sys.stderr)
            total = sum(graph[path[i]][path[i + 1]] for i in range(len(path) - 1))
            print(f"Длина пути: {total:.2f}", file=sys.stderr)


def prim_mst(graph, n, logger):
    mst = [[] for _ in range(n)]
    visited = [False] * n
    pq = []
    heappush(pq, (0, 0, -1))  # (вес, текущая_вершина, родитель)

    logger.log("Начинаем построение MST алгоритмом Прима...")

    while pq:
        weight, u, parent = heappop(pq)
        if visited[u]:
            continue

        visited[u] = True
        if parent != -1:
            mst[u].append((parent, weight))
            mst[parent].append((u, weight))
            logger.log(f"Добавляем ребро {parent}-{u} с весом {weight}")

        for v in range(n):
            if not visited[v] and graph[u][v] != -1:
                heappush(pq, (graph[u][v], v, u))

    logger.log_mst(mst)
    return mst


def tsp_2_approx(graph, start, n, logger):
    mst = prim_mst(graph, n, logger)

    visited = [False] * n
    path = []

    logger.log("\nНачинаем обход DFS для построения пути...")

    def dfs(u):
        visited[u] = True
        path.append(u)
        logger.log(f"Посещаем вершину {u}")

        for v, _ in sorted(mst[u], key=lambda x: x[1]):
            if not visited[v]:
                dfs(v)

    dfs(start)
    path.append(start)

    total_length = sum(graph[path[i]][path[i + 1]] for i in range(len(path) - 1))
    logger.log_path(path, graph)

    return total_length, path


def read_input():
    start = int(input("Введите стартовую вершину: "))
    print("Введите матрицу весов:")
    graph = []
    first_row = list(map(float, input().split()))
    graph.append(first_row)
    n = len(first_row)

    for _ in range(n - 1):
        row = list(map(float, input().split()))
        graph.append(row)

    return start, graph, n


def main():
    parser = ArgumentParser()
    parser.add_argument('-debug', action='store_true', help='Включить отладочный вывод')
    args = parser.parse_args()

    logger = DebugLogger(args.debug)
    start, graph, n = read_input()

    length, path = tsp_2_approx(graph, start, n, logger)

    print("Стоимость пути:", f"{length:.2f}")
    print("Путь:", " ".join(map(str, path)))

if __name__ == "__main__":
    main()