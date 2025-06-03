import heapq
import random
import sys
from argparse import ArgumentParser
from typing import List, Tuple, Optional


class DebugLogger:
    def __init__(self, debug: bool = False) -> None:
        self.debug = debug

    def log(self, message: str) -> None:
        if self.debug:
            print(f"[DEBUG] {message}", file=sys.stderr)

    def log_queue(self, pq: List[Tuple[float, float, int, List[int], List[bool]]]) -> None:
        if self.debug:
            print("\nТекущее состояние очереди приоритетов:", file=sys.stderr)
            for item in pq:
                bound, cost, curr, path, visited = item
                print(f"  bound={bound:.2f}, cost={cost:.2f}, vertex={curr}, path={path}", file=sys.stderr)

    def log_bound_calculation(self, curr_vertex: int, visited: List[bool], min_edges: List[List[float]], bound: float) -> None:
        if self.debug:
            print(f"\nРасчет нижней границы для вершины {curr_vertex}:", file=sys.stderr)
            print(f"Посещенные вершины: {visited}", file=sys.stderr)
            print("Минимальные ребра для каждой вершины:", file=sys.stderr)
            for i, edges in enumerate(min_edges):
                print(f"  Вершина {i}: {edges[0]:.2f} (и {edges[1]:.2f} если есть)", file=sys.stderr)
            print(f"Итоговая нижняя граница: {bound:.2f}", file=sys.stderr)

    def log_new_node(self, path: List[int], cost: float, bound: float, best_cost: float) -> None:
        if self.debug:
            print(f"\nОбработка узла:", file=sys.stderr)
            print(f"  Текущий путь: {' -> '.join(map(str, path))}", file=sys.stderr)
            print(f"  Текущая стоимость: {cost:.2f}", file=sys.stderr)
            print(f"  Нижняя граница: {bound:.2f}", file=sys.stderr)
            print(f"  Лучшая известная стоимость: {best_cost:.2f}", file=sys.stderr)

    def log_new_best(self, path: List[int], cost: float) -> None:
        if self.debug:
            print("\nНайден новый лучший путь!", file=sys.stderr)
            print(f"  Путь: {' -> '.join(map(str, path))}", file=sys.stderr)
            print(f"  Стоимость: {cost:.2f}", file=sys.stderr)

    def log_skip(self, vertex: int, new_cost: float, best_cost: float) -> None:
        if self.debug:
            print(f"  Пропуск вершины {vertex}: {new_cost:.2f} >= {best_cost:.2f}", file=sys.stderr)


def get_lower_bound(graph: List[List[float]], visited: List[bool], curr_vertex: int, n: int, logger: DebugLogger) -> float:
    if sum(visited) == n:
        return graph[curr_vertex][0]

    remaining = [v for v in range(n) if not visited[v]]
    bound = 0

    min_edges = []
    for v in range(n):
        edges = [graph[v][u] for u in range(n) if graph[v][u] != -1 and (not visited[u] or u == 0)]
        edges.sort()
        min_edges.append(edges[:2] if len(edges) >= 2 else [edges[0], float('inf')])

    bound += min_edges[curr_vertex][0]  # наименьшая стоимость ребра, которое можно использовать
                                        # для перехода из текущей вершины в следующую допустимую вершину

    for v in remaining:
        bound += min_edges[v][0] # добавляем минимальные рёбра оставшихся непосещённых вершин

    logger.log_bound_calculation(curr_vertex, visited, min_edges, bound)

    return bound


def tsp_branch_and_bound(n: int, graph: List[List[float]], logger: DebugLogger) -> Tuple[float, Optional[List[int]]]:
    visited = [False] * n
    visited[0] = True
    pq = [(0, 0, 0, [0], visited)]
    best_cost = float('inf')
    best_path = None

    logger.log(f"Начало алгоритма. Начальная вершина: 0")
    logger.log_queue(pq)

    while pq:
        bound, cost, curr, path, visited = heapq.heappop(pq)

        logger.log_new_node(path, cost, bound, best_cost)

        if len(path) == n:
            total_cost = cost + graph[curr][0]
            if total_cost < best_cost:
                best_cost = total_cost
                best_path = path[:]
                logger.log_new_best(best_path, best_cost)
            continue

        for next_vertex in range(n):
            if not visited[next_vertex]:
                new_cost = cost + graph[curr][next_vertex]
                if new_cost >= best_cost:
                    logger.log_skip(next_vertex, new_cost, best_cost)
                    continue

                new_visited = visited[:]
                new_visited[next_vertex] = True
                new_path = path + [next_vertex]
                lower_bound = get_lower_bound(graph, new_visited, next_vertex, n, logger)
                total_bound = new_cost + lower_bound

                if total_bound < best_cost:
                    heapq.heappush(pq, (total_bound, new_cost, next_vertex, new_path, new_visited))
                    logger.log(f"Добавлен новый узел в очередь: vertex={next_vertex}, path={new_path}, total_bound={total_bound:.2f}")
        if pq:
            logger.log_queue(pq)

    return best_cost, best_path

def generate_random_graph(n: int, min_cost: int = 1, max_cost: int = 100) -> list[list[float]]:
    graph = []
    for i in range(n):
        row = []
        for j in range(n):
            if i == j:
                row.append(-1)
            else:
                row.append(random.randint(min_cost, max_cost))
        graph.append(row)

    with open ("matrix.txt", "w") as file:
        for row in graph:
            file.write(" ".join(map(str, row)) + "\n")

    return graph


def main() -> None:
    parser = ArgumentParser()
    parser.add_argument('-debug', action='store_true', help='Включить отладочный вывод')
    parser.add_argument('-random', action='store_true', help='Сгенерировать случайную матрицу')
    args = parser.parse_args()

    n = int(input("Введите колличество городов: "))
    if args.random:
        graph = generate_random_graph(n)
    else:
        print("Введите матрицу стоимостей: ")
        graph = [list(map(float, input().split())) for _ in range(n)]

    logger = DebugLogger(debug=args.debug)

    total_cost, path = tsp_branch_and_bound(n, graph, logger)

    print("\nИтоговый результат:")
    print("Посещённые города:", " ".join(map(str, path)))
    print("Стоимость пути:", total_cost)


if __name__ == "__main__":
    main()