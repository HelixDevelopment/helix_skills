#!/usr/bin/env python3
"""Validate that CORPUS.yaml declares a well-formed, acyclic dependency graph.

Checks:
  1. Every requires/extends/recommends edge target resolves to a node
     declared in CORPUS.yaml (closed-world corpus; no dangling edges).
  2. The combined (requires + extends + recommends) directed graph contains
     no cycles.

Usage: python3 validate_dag.py [path/to/CORPUS.yaml]
Exit code 0 on success, 1 on failure (unresolved edges and/or cycle found).
"""
import sys
import yaml


def load_corpus(path):
    with open(path, "rb") as f:
        data = yaml.safe_load(f)
    return data["nodes"]


def build_graph(nodes):
    names = {n["name"] for n in nodes}
    edges = {}  # name -> list of (target, relation)
    unresolved = []
    for n in nodes:
        src = n["name"]
        edges[src] = []
        for relation in ("requires", "extends", "recommends"):
            for target in n.get(relation) or []:
                if target not in names:
                    unresolved.append((src, relation, target))
                edges[src].append((target, relation))
    return names, edges, unresolved


def find_cycle(names, edges):
    WHITE, GRAY, BLACK = 0, 1, 2
    color = {n: WHITE for n in names}
    path = []

    def visit(node):
        color[node] = GRAY
        path.append(node)
        for target, _relation in edges.get(node, []):
            if target not in color:
                continue  # unresolved edge, reported separately
            if color[target] == GRAY:
                cycle_start = path.index(target)
                return path[cycle_start:] + [target]
            if color[target] == WHITE:
                result = visit(target)
                if result:
                    return result
        path.pop()
        color[node] = BLACK
        return None

    for node in sorted(names):
        if color[node] == WHITE:
            cycle = visit(node)
            if cycle:
                return cycle
    return None


def main():
    path = sys.argv[1] if len(sys.argv) > 1 else "CORPUS.yaml"
    nodes = load_corpus(path)
    names, edges, unresolved = build_graph(nodes)

    edge_count = sum(len(v) for v in edges.values())

    if unresolved:
        print(f"UNRESOLVED EDGES ({len(unresolved)}):")
        for src, relation, target in unresolved:
            print(f"  {src} --{relation}--> {target}  [target not declared in CORPUS.yaml]")
        sys.exit(1)

    cycle = find_cycle(names, edges)
    if cycle:
        print("CYCLE FOUND:")
        print("  " + " -> ".join(cycle))
        sys.exit(1)

    print(f"DAG OK ({len(names)} nodes, {edge_count} edges)")
    sys.exit(0)


if __name__ == "__main__":
    main()
