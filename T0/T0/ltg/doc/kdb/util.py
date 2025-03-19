import os
import sys
import logging
import datetime as dt
from logging.handlers import RotatingFileHandler

pwd = os.path.dirname(os.path.realpath(__file__))

# 假设你有一个自定义的模块 my_module.py 存放在 /path/to/my/module 目录下，但是该目录并不在 Python 解释器的默认模块搜索路径中。你可以使用 sys.path.append() 方法将该目录添加到模块搜索路径中，从而使 Python 解释器能够找到并导入你的自定义模块
# sys.path.append(pwd+"/../")


today = dt.datetime.today().strftime('%Y%m%d')
log_format = logging.Formatter("%(asctime)s #-# %(name)s #-# %(levelname)s #-# Line:%(lineno)d #-# %(message)s \n")
log_handler = RotatingFileHandler(pwd + "/../log/ltg-{}.log".format(today), mode='a', maxBytes=50 * 1024 * 1024,
                                  backupCount=100, encoding=None, delay=0)
log_handler.setFormatter(log_format)
log_handler.setLevel(logging.INFO)

logger = logging.getLogger('root')
logger.addHandler(logging.StreamHandler(sys.stdout))
logger.addHandler(log_handler)
logger.setLevel(logging.INFO)
