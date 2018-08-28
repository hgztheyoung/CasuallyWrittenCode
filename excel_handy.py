# -*- coding: utf-8 -*-
#using code from https://www.jianshu.com/p/24a693fff4a3
#the original writer has all its copyright for that part of code

import xlrd
import csv

# read a excel sheet,and return an 2-dem array
# deal with merged slot
def read_sheet(file_path, sheet_index):
    # 打开文件
    workbook = xlrd.open_workbook(file_path)de
    sheet2 = workbook.sheet_by_index(sheet_index)  # sheet索引从0开始
    rows_num = sheet2.nrows
    cols_num = sheet2.ncols
    multilist = [['hugz002_dummy_init_should_not_appear_in_result' for col in range(cols_num)] for row in
                 range(rows_num)]
    for r in range(rows_num):
        # 一行数据的实体类
        for c in range(cols_num):
            cell_value = sheet2.row_values(r)[c]
            # print('第%d行第%d列的值：[%s]' % (r, c, sheet2.row_values(r)[c]))
            if (cell_value is None or cell_value == ''):
                cell_value = (get_merged_cells_value(sheet2, r, c))
                # 构建Entity
            multilist[r][c] = cell_value
    return multilist


def get_merged_cells(sheet):
    """
    获取所有的合并单元格，格式如下：
    """
    return sheet.merged_cells


def get_merged_cells_value(sheet, row_index, col_index):
    """
    newline =‘’
    如果是合并单元格，就返回合并单元格的内容
    :return:
    """
    merged = get_merged_cells(sheet)
    for (rlow, rhigh, clow, chigh) in merged:
        if (row_index >= rlow and row_index < rhigh):
            if (col_index >= clow and col_index < chigh):
                cell_value = sheet.cell_value(rlow, clow)
                # print('该单元格[%d,%d]属于合并单元格，值为[%s]' % (row_index, col_index, cell_value))
                return cell_value
                break
    return None


def genUniqueRows(twodimArray):
    if len(twodimArray) <= 1:
        return twodimArray
    ret = [twodimArray[0]]
    for i in range(1, len(twodimArray)):
        if twodimArray[i] != twodimArray[i - 1]:
            ret.append(twodimArray[i])
    return ret

def csvConvert(from_csv_filepath,res_file_path,
               from_cols_list,to_col_list):
    '''
        2 list should be the same length
        use a csv containing from_cols_list to
        generate result csv with header to_col_list
    '''
    if len(from_cols_list)!= len(to_col_list):
        return
    with open(from_csv_filepath) as fromcsvfile:
        reader = csv.DictReader(fromcsvfile)
        with open(res_file_path,'w',newline ='') as rescsvfile:
            writer = csv.DictWriter(rescsvfile, fieldnames=to_col_list)
            writer.writeheader()
            for row in reader:
                d = {to_col_list[i]:row[from_cols_list[i]]
                     for i in range(0,len(to_col_list))}
                writer.writerow(d)

def writeTwodimlistToCsv(csvpath):
    with open(csvpath, 'w', newline='') as csvfile:
        spamwriter = csv.writer(csvfile, delimiter=',',
                                quotechar='|', quoting=csv.QUOTE_MINIMAL)
        for row in res:
            spamwriter.writerow(row)

if __name__ == "__main__":
    rawres = read_sheet('erqi_data.xlsx', 0)
    res = genUniqueRows(rawres)
    writeTwodimlistToCsv('eggs.csv')
    csvConvert('eggs.csv','neweggs.csv',['a','b','d'],['x','y','z'])
