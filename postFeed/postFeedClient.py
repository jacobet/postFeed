
import sys
import argparse
import json
import requests
import socket
from datetime import datetime

def arg_list() ->  list:
    parser = argparse.ArgumentParser()
    parser.add_argument('cmd', type=str, help='Command option: post, get-post, search-posts, delete-post, edit-post, like, dislike, order')
    parser.add_argument('--post-id', required=False, type=int, help='Post\'s id')
    parser.add_argument('--author', required=False, type=str, help='Post\'s Author')
    parser.add_argument('--content', required=False, type=str, help='Post\'s Content')
    parser.add_argument('--create-at', required=False, 
                        type=lambda s: datetime.strptime(s, '%d-%m-%Y'), 
                        help='Posts\' created date that bigger than param in format \'dd-mm-yyyy\'')
    parser.add_argument('--update-at', required=False, 
                        type=lambda s: datetime.strptime(s, '%d-%m-%Y'), 
                        help='Posts\' updated date that bigger than param in format \'dd-mm-yyyy\'')
    parser.add_argument('--like', required=False, type=int, help='Post\'s Likes')
    parser.add_argument('--dislike', required=False, type=int, help='Post\'s Dislikes')
    parser.add_argument('--order', required=False, type=str, help='Post\'s Order')
    parser.add_argument('--skip', required=False, type=int, help='Post\'s Skip')
    parser.add_argument('--limit', required=False, type=int, help='Post\'s Limit')
    parser.add_argument('--page', required=False, type=int, help='Post\'s Page')
    
    args = parser.parse_args()
    cmd = args.cmd
    post_id = args.post_id
    author = args.author
    content = args.content
    create_at = args.create_at
    update_at = args.update_at
    like = args.like
    dislike = args.dislike
    order = args.order
    skip = args.skip
    limit = args.limit
    page = args.page
    return [cmd, post_id, author, content, create_at, update_at, like, dislike, order, skip, limit, page]

def get_params(post_id, author, content, create_at = None, update_at = None, 
               like = None, dislike = None, order = None, skip = None, limit = None, page = None) ->  dict:
    data = {}
    if author is not None:
        data['author'] = author
    if content is not None:
        data['content'] = content
    if create_at is not None:
        data['create_at'] = create_at
    if update_at is not None:
        data['update_at'] = update_at
    if like is not None:
        data['like'] = like
    if dislike is not None:
        data['dislike'] = dislike
    if order is not None:
        data['order'] = order
    if skip is not None:
        data['skip'] = skip
    if limit is not None:
        data['limit'] = limit
    if page is not None:
        data['page'] = page
        
    return data
    
def is_open(ip: str, port: int) ->  bool:
   s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
   try:
      s.connect((ip, port))
      return True
   except:
      return False

def post_fun(address: str) ->  requests.models.Response:
    
    my_list = arg_list()
    cmd = my_list[0]
    post_id = my_list[1]
    author = my_list[2]
    content = my_list[3]
    create_at = my_list[4]
    update_at = my_list[5]
    like = my_list[6]
    dislike = my_list[7]
    order = my_list[8]
    skip = my_list[9]
    limit = my_list[10]
    page = my_list[11]

    req = ''
    address += 'posts'
    if cmd == 'post':
        if author is not None and content is not None:
            data = {'author' : author,
                    'content' : content}
            req = requests.post(address, json=data)
        else:
            print( 'Enter author and content arguments' )
    elif cmd == 'get-post':
        if post_id is not None:
            address += '/{0}'.format( post_id )
            req = requests.get(address)
        else:
            print( 'Enter post-id argument' )
    elif cmd == 'search-posts':
            data = get_params(post_id, author, content, create_at, update_at, like, dislike, order, skip, limit, page)
            if len(data) > 0:
                req = requests.get(address, json=data)
            else:
                req = requests.get(address)
    elif cmd == 'delete-post':
        if post_id is None:
            print( 'Enter post-id argument' )
        else:
            address += '/{0}'.format( post_id )
            req = requests.delete(address)
    elif cmd == 'edit-post' or cmd == 'like' or cmd == 'dislike':
        if post_id is not None:
            if cmd == 'edit-post':
                data = get_params(post_id, author, content, create_at, update_at)
            elif cmd == 'like':
                data = {'post_id' : post_id,
                        'like' : 1}
            elif cmd == 'dislike':
                data = {'post_id' : post_id,
                        'dislike' : 1}
            address += '/{0}'.format( post_id )
            req = requests.put(address, json=data)
        else:
            print( 'Enter post-id argument' )
        
    if req != '' and req is not None:
        return req
    else:
        return None

def comment_fun(address: str) ->  requests.models.Response:
    
    my_list = arg_list()
    cmd = my_list[0]
    post_id = my_list[1]
    author = my_list[2]
    content = my_list[3]
    create_at = my_list[4]
    update_at = my_list[5]
    like = my_list[6]
    dislike = my_list[7]
    order = my_list[8]
    skip = my_list[9]
    limit = my_list[10]
    page = my_list[11]

    
    address += 'comments'
    if cmd == 'comment':
        if post_id is not None and author is not None and content is not None:
            data = {'postID' : post_id,
                    'author' : author,
                    'content' : content}
            req = requests.post(address, json=data)
            return req
        else:
            print( 'Enter correct command' )
            
    return None

def print_response(req: requests.models.Response):
    content = req.content.decode("utf-8")
    if req.status_code < 200 or req.status_code >= 300:
        result = 'Error! status_code - {0}, post-id {1} {2}'.format( req.status_code, content.strip(), req.reason )
        print(result)
    else:
        parsed = json.loads(content)
        if len(parsed) == 1:
            parsed = parsed[0]
        prettyJson = json.dumps(parsed, indent=4)
        print(prettyJson)

if __name__ == "__main__":
    localhost = '127.0.0.1'
    port = 1323
    conn = is_open(localhost, port)
    if not conn:
        print( 'Server is shutdown. Please Start it and play again' )
        sys.exit()
    
    
    address = 'http://{0}:{1}/'.format( localhost, port )
  
    req1 = post_fun(address)    
    req2 = comment_fun(address)  
    if req1 is not None:
        print_response(req1)
    elif req2 is not None:
        print_response(req2)
    else:
        print( 'Enter correct command' )
        print( 'Command option: post, commant, get-post, search-posts, delete-post, edit-post, like, dislike, order' )
        sys.exit()
    
    

