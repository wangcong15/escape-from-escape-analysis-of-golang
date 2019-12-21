# -*- coding: utf-8 -*-

# =============================== #
# Author : Cong Wang              #
# Email  : bryantwangcong@163.com #
# =============================== #

import os
import re
import logging
import argparse
import subprocess
from common import findGoSrc, setDir

LOGGER_LEVEL = 10 # 10: debug, 20: info

# Overall Handler for the task
class Go2IrHandler(object):


    def __init__(self, args):
        # basic member variables
        self.args = args
        self.project = args.project
        self.output = args.output if args.output != "" else os.path.join(os.getcwd(), "../go2ir_output")
        self.logger = logging.getLogger("default")
        self.error_status = False
        self.unique_ll_id = 1
        self.commit_ids = []
        self.commit_update_files = []
        self.curr_commit_idx = 0
        self.modified_files = []

        # basic initialization
        self.logger_init()
        self.output_init()
        self.check_args()


    # logger_init : logger initialization
    def logger_init(self):
        # logger level and format
        self.logger.setLevel(LOGGER_LEVEL)
        log_format = logging.Formatter('%(asctime)s - %(name)s - %(levelname)s - %(message)s')

        # declare the console logger
        console_handler = logging.StreamHandler()
        console_handler.setLevel(LOGGER_LEVEL)
        console_handler.setFormatter(log_format)
        self.logger.addHandler(console_handler)


    # output_init : output initialization
    def output_init(self):
        setDir(self.output)


    # check_args : check the path of ${project}
    def check_args(self):
        if not(os.path.exists(self.project) and os.path.isdir(self.project)):
            self.logger.error("Path of ${project} is invalid: %s" % (self.project))
            self.error_status = True
        else:
            self.logger.debug("Path of ${project} is valid: %s" % (self.project))


    # run : execution 
    def run(self):
        self.logger.debug(self.args)
        setDir(self.output)
        self.handle_dir(os.path.abspath(self.project))
        

    # handle_dir : handler directory recursively
    def handle_dir(self, dir):
        files = os.listdir(dir)
        file_paths = [os.path.join(dir, file) for file in files]
        
        if findGoSrc(file_paths):
            # if we find go source files, this part is significant!
            build_script = "go build -work -x *.go 1> transcript.txt 2>&1"
            self.logger.debug("%s ===> %s" % (build_script, dir))
            subprocess.Popen(build_script, cwd=dir, shell=True).wait()

            transcript_path = os.path.join(dir, "transcript.txt")
            new_script = []
            with open(transcript_path) as f:
                trans_text = f.read().split("\n")
                if trans_text[0].startswith('WORK='):
                    new_script.append(trans_text[0])
                pattern_cd = "cd %s" % dir
                for idx, text in enumerate(trans_text):
                    if text == pattern_cd and trans_text[idx+1].find("llvm-goc") >= 0:
                        modified_script = trans_text[idx+1] + " -S -emit-llvm"
                        modified_script = modified_script.replace("-o $WORK/b001/_go_.o", "-o %s/%s.ll" % (self.output, self.unique_ll_id))
                        new_script.append(modified_script)
                        break

            if len(new_script) == 2:
                self.unique_ll_id += 1
                generate_script = " && ".join(new_script)
                self.logger.debug(generate_script)
                subprocess.Popen(generate_script, cwd=dir, shell=True).wait()
            else:
                self.logger.error(new_script)

        # handle sub-dirs recursively
        # for file in file_paths:
        #     if os.path.exists(file) and os.path.isdir(file):
        #         dir_name = file[file.rfind("/")+1:]
        #         if dir_name != "vendor":
        #             self.handle_dir(file) 


def main():
    # Define command line arguments
    arg_parser = argparse.ArgumentParser()
    arg_parser.add_argument('-p', '--project', help='choose golang project', default='')
    arg_parser.add_argument('-o', '--output', help='choose output dir', default='')
    args = arg_parser.parse_args()
    
    handler = Go2IrHandler(args)
    handler.run()

if __name__ == "__main__":
    main()