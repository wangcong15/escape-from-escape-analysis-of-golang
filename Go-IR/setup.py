from setuptools import setup  
  
setup(  
    name = "Go2IR",  
    version = "0.0.1",  
    keywords = ("golang", "llvm ir"),
    description = "Go2ir is a tool for converting a go language program project into LLVM IR code based on the \"gollvm\" tool",  
    long_description = "",  
    license = "Tsinghua Univ. Licence",  
  
    url = "",  
    author = "Cong Wang, Yu Jiang, Jian Gao",  
    author_email = "wangcong15@mails.tsinghua.edu.cn",  

    packages = ['go2ir'],  
    package_dir = {'': 'src'},
    package_data = {},
    include_package_data = False,  
    platforms = "any",  
    install_requires = [],  

    scripts = [],  
    entry_points = {  
        'console_scripts': [  
            'Go2IR=go2ir.main:main'  
        ]  
    }  
)