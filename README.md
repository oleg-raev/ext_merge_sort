# Внешняя сортировка слиянием
_Для сортировки фрагментов используется метод сортировки выбором_

## Допущения:
1. Сортировка расчитана исключительно на однобайтовые кодировки (не сложно поддержать и многобайтовые, но отвлекаться
на это не хотелось)
2. Считаем, что порядок символом должен соответствовать порядку их расположения в таблице ASCII
3. В условиях временных ограничений не ставилось целью сделать хорошее оформление кода и разбиение на структуры
4. Решение не претендует быть оптимальным и требует доработки (как архитектурной, так и логической)
5. Для удобства визуальной проверки алгоритма, генерация файла использует только символы от A до Z

## Выбранный метод решения
1. Исходный файл разбивается на фрагменты по 1000 строк
2. Каждый фрагмент сортируется методом сортировки выбора и сохраняется во временную папку
3. Отсортированные фрагменты объединяются методом сортировки слиянием до тех пор, пока не останется одного файла

## Как собрать и проверить
1. Установить [Go 1.13.*](https://golang.org/dl/)
2. Склонировать репозиторий ```git clone https://github.com/oleg-raev/ext_merge_sort.git```
3. Зайти в папку репозитория в терминале `cd ext_merge_sort`
4. Собрать проект `GO111MODULE=on go build` (если не удается, попробуйте загрузить зависимости `GO111MODULE=on go mod tidy`)
5. Исполняемый файл будет располагаться по адресу ./ext_merge_sort

## Работа с приложением
### Помощь по параметрам
```ext_merge_sort -help```

### Сгенерировать файл с данными
```
ext_merge_sort -generate \
  -out="{адрес файла назначения}" \
  -lines={количество строк} \
  -rowlen={максимальная длина строки}
```

### Сортировать файл
```
ext_merge_sort -sort \
  -in="{адрес файла с данными}"
  -out="{адрес файла назначения}"
```
