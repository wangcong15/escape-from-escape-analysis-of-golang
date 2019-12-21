# -*- coding: utf-8 -*-

# =============================== #
# Author : Cong Wang              #
# Email  : bryantwangcong@163.com #
# =============================== #

import os
import shutil

# create directory, remove if exist
def setDir(dir_path):
    if os.path.exists(dir_path):
        shutil.rmtree(dir_path)
    os.mkdir(dir_path)


# findGoSrc : return True if find go source file in ${file_list}
def findGoSrc(file_list):
    for file in file_list:
        if os.path.exists(file) and os.path.isfile(file) and file.endswith(".go"):
            return True
    return False
