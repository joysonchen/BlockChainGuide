# 准备工作：
#   pip install  pipenv
#   pipenv --python=python3.6.5 或 pipenv --python3.6.5
#   ls 生成pipfile文件
#   cat Pipfile 查看pipfile
#   pipenv install flask==0.12.2
#   pipenv install requests==2.18.4
#
# json结构：{
#     "index":0,
#     "timestamp":"",
#     "transpations":[
#         {
#             "sender":"",
#             "recipient":"",
#             "amount":""
#         }
#     ],
#     "proof":"",
#     "pervious_hash":""
# }
import hashlib
import json
import time
from argparse import ArgumentParser
from urllib.parse import urlparse

from uuid import uuid4

# from flask import Flask, jsonify, request
from flask import Flask, jsonify, request
from pipenv.vendor import requests


class Blockchain:
    # 构造函数
    def __init__(self):
        self.chain = []
        self.current_transcations = []
        self.nodes = set()
        # 创世区块
        self.new_block(proof=100, pervious_hash=1)

    # 注册节点,如http://192.168.1.1  http://192.168.1.12
    def register_node(self, address: str):
        parsed_url = urlparse(address)
        self.nodes.add(parsed_url.netloc)

    # 验证链合法性,每个块的hash是否是上个块的hash
    def valid_chain(self, chain) -> bool:
        last_block = chain[0]
        current_index = 1
        while current_index < len(chain):
            block = chain[current_index]
            if (block['pervious_hash'] != self.hash(last_block)):
                return False

            if not self.valid_proof(last_block['proof'], block['proof']):
                return False

            last_block = block
            current_index += 1

        return True

    #  解决节点冲突，解决分叉，达成共识：谁的链长，取谁
    def resolve_conflicts(self) -> bool:
        neighbours = self.nodes
        max_lenght = len(self.chain)
        for node in neighbours:
            response = requests.get(f'http://{node}/chain')
            new_chain = None

            if response.status_code == 200:
                length = response.json()['length']
                chain = response.json()['chain']

                if length > max_lenght and self.valid_chain(chain):
                    max_lenght = length
                    new_chain = chain

            if new_chain:
                self.chain = new_chain
                return True
            return False


    def new_block(self, proof, pervious_hash=None):
        block = {
            'index': len(self.chain) + 1,
            'timestamp': time.time(),
            'transcations': self.current_transcations,
            'proof': proof,
            "pervious_hash": pervious_hash or self.hash(self.last_block)
        }
        self.current_transcations = []
        self.chain.append(block)
        return block

    # 添加新交易
    def new_transaction(self, sender, recipient, amount) -> int:
        self.current_transcations.append({
            'sender': sender,
            'recipient': recipient,
            'amount': amount
        })
        return self.last_block['index'] + 1

    # hash算法
    @staticmethod
    def hash(block):
        block_string = json.dumps(block, sort_keys=True).encode()
        return hashlib.sha256(block_string).hexdigest()

    # 获取最后区块
    @property
    def last_block(self):
        return self.chain[-1]

    # pow 验证
    def proof_of_work(self, last_proof: int) -> int:
        proof = 0
        while self.valid_proof(last_proof, proof) is False:
            proof += 1
            print("工作量证明：" + str(proof))
        return proof

    # 权益验证
    @staticmethod
    def valid_proof(last_proof: int, proof: int) -> bool:
        guess = f'{last_proof}{proof}'.encode()
        guess_hash = hashlib.sha256(guess).hexdigest()
        print("验证证明：" + str(guess_hash))
        return guess_hash[0:4] == "0000"


app = Flask(__name__)
node_identifier = str(uuid4()).replace('-', '')
blockchain = Blockchain()

# pow机制，挖矿
@app.route('/mine', methods=['GET'])
def mine():
    last_block = blockchain.last_block
    last_proof = last_block['proof']
    proof = blockchain.proof_of_work(last_proof)
    # 给自己奖励
    blockchain.new_transaction(sender="0", recipient=node_identifier, amount=1)
    block = blockchain.new_block(proof, None)
    response = {
        "message": "New Block Forged",
        "index": block['index'],
        "transcations": block['transcations'],
        "proof": block['proof'],
        "pervious_hash": block['pervious_hash']
    }
    return jsonify(response), 200

# 请求链信息
@app.route('/chain', methods=['GET'])
def full_chain():
    response = {
        'chain': blockchain.chain,
        'length': len(blockchain.chain)
    }
    return jsonify(response), 200


# {"nodes":["http://127.0.0.1:5000"]}
# 注册节点
@app.route('/nodes/register', methods=['POST'])
def register_nodes():
    values = request.get_json()
    nodes = values.get("nodes")
    if nodes is None:
        return "Error : please supply a valid list of nodes", 400
    for node in nodes:
        blockchain.register_node(node)
    response = {
        "message": "new nodes have been added",
        "total_nodes": list(blockchain.nodes)
    }
    return jsonify(response), 201

# 创建新交易
@app.route('/transcations/new', methods=['POST'])
def new_transaction():
    values = request.get_json()
    required = ["sender", "recipient", "amount"]
    if values is None:
        return "missing values", 400
    if not all(k in values for k in required):
        return "missing values", 400
    index = blockchain.new_transaction(values['sender'], values['recipient'], values['amount'])
    response = {"message": f'Transcation will be added to Block{index}'}
    return jsonify(response), 201


# 解决冲突，达到共识机制
@app.route('/nodes/resolve', methods=['GET'])
def consensus():
    replaced = blockchain.resolve_conflicts()
    if replaced:
        response = {
            'message': 'Our chain was replaced',
            'new_chain': blockchain.chain
        }
    else:
        response = {
            'message': 'Our chain is authoritative',
            'chain': blockchain.chain
        }
    return jsonify(response), 200

# 动态注册多个端口
if __name__ == '__main__':
    parser=ArgumentParser()
    # -p --port 5001
    parser.add_argument('-p','--port',default=5000,type=int,help='port need to listen ')
    args=parser.parse_args()
    port = args.port
    app.run(host='0.0.0.0', port=port)
