import os
import json
import requests

from common.ripple import userUtils

from . import flavours

from .discord_hooks import Webhook

def call(m, *args, user_id = None, discord_m = False, embed_args = None):
    try:
        if flavours.config is None:
            cache_config()

        username = None
        if user_id is not None:
            username = userUtils.getUsername(user_id)
            m = m.replace("USERNAME()", username)

        if flavours.config["webhook"]["enable"] and discord_m:
            if embed_args is not None:
                embed = Webhook(flavours.config["webhook"]["url"], 
                                **embed_args,
                                footer="Caker by Sunpy @osufx",
                                footer_icon="http://i.imgur.com/NCYspz8.png"
                            )
            else:
                embed = Webhook(flavours.config["webhook"]["url"], 
                                msg=m
                            )

            embed.post()
    except Exception as e:
        s_print("Unable to call police; {}".format(str(e)))
    
    s_print(m)

def cache_config():
    with open(os.path.join(os.path.dirname(__file__), "config.json"), "r") as f:
        flavours.config = json.load(f)
    s_print("Config was loaded. We are ready to go!")

def s_print(m):
    print("[Police] {}".format(m))