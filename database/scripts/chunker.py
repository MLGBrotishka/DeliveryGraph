"""Чанкер для проекта Доставки группы М8О-206Б-21

Этот скрипт обрабатывает графовое предсталвение участка карты OpenStreetMap.
Получая на вход исходный файл-граф (SOURCE_FILE), скрипт разбивает граф на
указанное пользователем количество чанков (CHUNK_NUM^2).

Чанки нумеруются в декартовой системе координат, начиная с чанка (0,0).

|     |     |
|     |     |
|_ _ _|_ _ _|_ _ _
|     |     |
| 0,1 | 1,1 |
|_ _ _|_ _ _|_ _ _
|     |     |
| 0,0 | 1,0 |
|_ _ _|_ _ _|_ _ _

Информация о каждом чанке сохраняется в двух файлах:
    V_x_y.txt : перечисление вершин и их координат
    E_x_y.txt : перечисление дуг и времени(сек),
        необходимого на их преодоление

Написал(Заговнокодил): Медведев Кирилл (МК)
"""


import os
from math import radians, sin, cos, sqrt, atan2


EPS = 0.05
SOURCE_FILE = 'ma.pypgr'


def create_grid(chunk_num: str, source: str) -> tuple:
    
    with open(source, 'r') as source_file:

        max_lat = float('-inf')
        min_lat = float('inf')
        max_long = float('-inf')
        min_long = float('inf')
        
        for i, line in enumerate(source_file):
            numbers = line.split()
            
            if i >= 9 and len(numbers) == 3:
                # Находим минимальные и максимальные значения широты и долготы
                numbers = line.split()
                if (float(numbers[2]) > max_long):
                    max_long = float(numbers[2])
                if (float(numbers[2]) < min_long):
                    min_long = float(numbers[2])
                if (float(numbers[1]) > max_lat):
                    max_lat = float(numbers[1])
                if (float(numbers[1]) < min_lat):
                    min_lat = float(numbers[1])

        # Отступим от самых крайних вершин ВСЕГО графа на EPS
        min_lat, max_lat = min_lat + EPS, max_lat - EPS
        min_long, max_long = min_long + EPS, max_long - EPS

        print("min_lat:", min_lat)
        print("max_lat:", max_lat)
        print("min_long:", min_long)
        print("max_long:", max_long)
        
        # Получим размер одного чанка
        latitude_step = (max_lat - min_lat) / chunk_num
        longtitude_step = (max_long - min_long) / chunk_num
        print("chunk_size:", latitude_step, longtitude_step)

        # Получим границы чанков
        horizontal_grid = []
        vertical_grid = []
        for i in range(chunk_num + 1):
            horizontal_grid.append(min_lat + latitude_step * i)
            vertical_grid.append(min_long + longtitude_step * i)
        
        return horizontal_grid, vertical_grid


def get_max_node_num(source_file: str) -> int:
    """Функция, возвращающая максимальное значение номера вершины в рассматриваемом графе."""
    
    max_node_num = int(-1)
    with open(source_file, 'r') as source_file:
        for i, line in enumerate(source_file):
            numbers = line.split()
            
            if i >= 9 and len(numbers) == 3:
                if (int(numbers[0]) > max_node_num):
                    max_node_num = int(numbers[0])

    return max_node_num


def intersection_point(lat1: float, lon1: float, lat2: float, lon2: float, lat: float=None, lon: float=None) -> tuple:
    if lat is not None:
        lon = lon1 + (lat - lat1) * (lon2 - lon1) / (lat2 - lat1)
    else:
        lat = lat1 + (lon - lon1) * (lat2 - lat1) / (lon2 - lon1)

    return lat, lon


def distance(lat1:float, lon1:float, lat2:float, lon2:float) -> float:
    # Конвертируем координаты в радианы
    lat1, lon1, lat2, lon2 = map(radians, [lat1, lon1, lat2, lon2])

    # Вычисляем разницу между координатами
    dlat = lat2 - lat1
    dlon = lon2 - lon1

    # Вычисляем расстояние по формуле гаверсинуса
    a = sin(dlat/2)**2 + cos(lat1) * cos(lat2) * sin(dlon/2)**2
    c = 2 * atan2(sqrt(a), sqrt(1-a))
    R = 6371    # Радиус Земли в километрах
    res = R * c * 1000.0
    
    return float(round(res, 2))


def break_intersected_edges(horizontal_grid: list, vertical_grid: list, source: str, output: str) -> int:

    new_edges_num = int(0)
    current_node_number = get_max_node_num(source) + 1

    nodes_to_tmp = set()
    
    with open(source, 'r') as source_file, open(output, 'w') as destination_file:

        all_lines = source_file.readlines()
        source_file.seek(0)
        
        for i, line in enumerate(source_file):
            numbers = line.split()
            
            if i < 9:
                destination_file.write(line)
                
            elif len(numbers) == 3:
                pass
                destination_file.write(line)
                
            elif len(numbers) == 6:
                new_nodes = set() # Node: (latitude, longtitude)
                
                start_node_num = int(numbers[0])
                end_node_num = int(numbers[1])
                       
                end_node_inf = all_lines[end_node_num + 9].split()
                end_node_lat = float(end_node_inf[1])
                end_node_long = float(end_node_inf[2])
                
                start_node_inf = all_lines[start_node_num + 9].split()
                start_node_lat = float(start_node_inf[1])
                start_node_long = float(start_node_inf[2])
                
                if (start_node_long > end_node_long):
                    start_node_long, end_node_long = end_node_long, start_node_long
                    start_node_lat, end_node_lat = end_node_lat, start_node_lat
                    start_node_num, end_node_num = end_node_num, start_node_num

                
                for grid_line in vertical_grid:
                    if (start_node_long < grid_line) and (grid_line < end_node_long):
                        new_node_coord = intersection_point(
                                                            start_node_lat, start_node_long,
                                                            end_node_lat, end_node_long,
                                                            lon=grid_line)

                        new_nodes.add((new_node_coord))
                    

                for grid_line in horizontal_grid:
                    if (start_node_lat < grid_line) and (grid_line < end_node_lat):
                        new_node_coord = intersection_point(
                                                            start_node_lat, start_node_long,
                                                            end_node_lat, end_node_long,
                                                            lat=grid_line)

                        new_nodes.add((new_node_coord))


                if len(new_nodes) == 0:
                    destination_file.write(line)
                else:
                    new_nodes.add((start_node_lat, start_node_long))
                    new_nodes.add((end_node_lat, end_node_long))
                    new_nodes = sorted(new_nodes)

                    
                    for edge_ind in range(len(new_nodes) - 1):
                        current_edge_line = []
                        
                        if edge_ind == 0:
                            current_edge_line.append(str(start_node_num))
                            current_edge_line.append(str(current_node_number))

                            nd = str(current_node_number) + " " + str(new_nodes[edge_ind + 1][0]) + " " + str(new_nodes[edge_ind + 1][1]) + "\n"
                            nodes_to_tmp.add(nd)
                            
                
                        elif edge_ind == (len(new_nodes) - 2):
                            current_edge_line.append(str(current_node_number - 1))
                            current_edge_line.append(str(end_node_num))

                            
                        else:
                            current_edge_line.append(str(current_node_number - 1))
                            current_edge_line.append(str(current_node_number))

                            nd = str(current_node_number) + " " + str(new_nodes[edge_ind + 1][0]) + " " + str(new_nodes[edge_ind + 1][1]) + "\n"
                            nodes_to_tmp.add(nd)


                        current_node_number += 1 
                        current_edge_line.append(str(distance(new_nodes[edge_ind][0], new_nodes[edge_ind][1],
                                                              new_nodes[edge_ind+1][0], new_nodes[edge_ind+1][1])))
                        current_edge_line.append(str(numbers[3]))
                        current_edge_line.append(str(numbers[4]))
                        current_edge_line.append(str(numbers[5]))

                        destination_file.write( (str(" ".join(current_edge_line)) + str("\n")) )
                    
                        
                    new_edges_num += len(new_nodes) - 1

    with open(output, 'r') as file:
        lines = file.readlines()
    lines[7] = str(int(lines[7]) + len(nodes_to_tmp)) + str("\n")
    lines[8] = str(int(lines[8]) + new_edges_num) + str("\n")
    with open(output, 'w') as file:
        file.writelines(lines)

    with open(output, 'r') as file:
        lines = file.readlines()
        lines = lines[:9] + sorted(list(nodes_to_tmp)) + lines[9:]
    with open(output, 'w') as file:
        file.writelines(lines)

    return len(nodes_to_tmp)
                        

def get_chunk(chunk_num: int, grid_x: int, grid_y: int, source: str, folder_name: str):

    chunk_num = int(chunk_num)
    grid_x = int(grid_x)
    grid_y = int(grid_y)
    number_of_selected_edges = int(0)

    # Создаем файл-представление графа в чанке
    vertices_filename = f'V_{grid_x}_{grid_y}.txt'
    edges_filename = f'E_{grid_x}_{grid_y}.txt'
    vertices_file_path = os.path.join(folder_name, vertices_filename)
    edges_file_path = os.path.join(folder_name, edges_filename)
    
    with open(source, 'r') as source_file, open(vertices_file_path, 'w') as verts, open(edges_file_path, 'w') as edgs:

        max_lat = float('-inf')
        min_lat = float('inf')
        max_long = float('-inf')
        min_long = float('inf')
        
        for i, line in enumerate(source_file):
            numbers = line.split()
            if i < 9:
                # Заполняем первые 9 служебных строк в файле-представлении графа
                pass
                #destination_file.write(line)
            elif len(numbers) == 3:
                # Находим минимальные и максимальные значения широты и долготы
                numbers = line.split()
                if (float(numbers[2]) > max_long):
                    max_long = float(numbers[2])
                if (float(numbers[2]) < min_long):
                    min_long = float(numbers[2])
                if (float(numbers[1]) > max_lat):
                    max_lat = float(numbers[1])
                if (float(numbers[1]) < min_lat):
                    min_lat = float(numbers[1])

        # Отступим от самых крайних вершин ВСЕГО графа на EPS
        min_lat, max_lat = min_lat + EPS, max_lat - EPS
        min_long, max_long = min_long + EPS, max_long - EPS

        # Получим размер одного чанка
        gor_step = (max_lat - min_lat) / chunk_num
        vert_step = (max_long - min_long) / chunk_num

        # Получим границы чанка
        gor_grid_1 = min_lat + grid_y * gor_step
        gor_grid_2 = gor_grid_1 + gor_step
        
        vert_grid_1 = min_long + grid_x * vert_step
        vert_grid_2 = vert_grid_1 + vert_step

        # Вернем каретку в начало файла-источника, после поиска минимальных величин. 
        source_file.seek(0)
        
        number_of_selected_edges = int(0)
        selected_nodes = set()

        for i, line in enumerate(source_file):
            numbers = line.split()
            
            if i >= 9:
                # Получим все вершины, входящие в чанк
                if len(numbers) == 3:
                    lat = float(numbers[1])
                    long = float(numbers[2])
                    node_number = int(numbers[0])
                
                    if (lat >= gor_grid_1
                            and lat <= gor_grid_2
                            and long >= vert_grid_1
                            and long <= vert_grid_2):

                        selected_nodes.add(node_number)

                        new_line = []
                        for i in range(3):
                            new_line.append(str(numbers[i]))

                        new_line.append(str(grid_x))
                        new_line.append(str(grid_y))
                        
                        if (lat == gor_grid_1):
                            new_line.append(str(grid_x))
                            new_line.append(str(grid_y - 1))
                        if (lat == gor_grid_2):
                            new_line.append(str(grid_x))
                            new_line.append(str(grid_y + 1))
                        if (long == vert_grid_1):
                            new_line.append(str(grid_x - 1))
                            new_line.append(str(grid_y))
                        if (long == vert_grid_2):
                            new_line.append(str(grid_x + 1))
                            new_line.append(str(grid_y))
                            
                        
                        verts.write(" ".join(new_line) + str("\n"))

                
                            
                        
                # Получим все дуги, входящие в чанк
                elif len(numbers) == 6:
                    start_node = int(numbers[0])
                    end_node = int(numbers[1])

                    if (start_node in selected_nodes) and (end_node in selected_nodes):
                        
                        new_line = []
                        new_line.append(str(start_node))
                        new_line.append(str(end_node))
                        new_line.append( str(round(float(numbers[2]) / (float(numbers[4]) / 3.6), 4)) )
                        edgs.write(" ".join(new_line) + str("\n"))
                        
                        if (int(numbers[5]) == 1):
                            new_line[0], new_line[1] = new_line[1], new_line[0]
                            edgs.write(" ".join(new_line) + str("\n"))
                            number_of_selected_edges += 1
                        
                        number_of_selected_edges += 1
                        

                
    # Укажем в служебном блоке файла-представления графа число вершин и дуг
    with open(vertices_file_path, 'r') as f1, open(edges_file_path, 'r') as f2:
        v = f1.readlines()
        e = f2.readlines()
    
    with open(vertices_file_path, 'w') as f1, open(edges_file_path, 'w') as f2:
        v = [str(str(len(selected_nodes)) + str("\n"))] + v
        e = [str(str(number_of_selected_edges) + str("\n"))] + e
        f1.writelines(v)
        f2.writelines(e)

        
        
if __name__ == "__main__":

    CHUNK_NUM = 5
    
    os.makedirs('chunks', exist_ok=True)

    hor_grid, vert_grid = create_grid(CHUNK_NUM, SOURCE_FILE)
        
    break_intersected_edges(hor_grid, vert_grid, SOURCE_FILE, 'tmp.txt')

    for i in range(CHUNK_NUM):
        for j in range(CHUNK_NUM):
            get_chunk(CHUNK_NUM, i, j, 'tmp.txt', 'chunks')
